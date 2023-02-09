package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"unicode"
)

const (
	MIN_NUM_OF_TESTS = 1
	MAX_NUM_OF_TESTS = 100

	MIN_NUM_OF_CARS = 0
	MAX_NUM_OF_CARS = 500

	MIN_NUM_OF_EVENTS = 0
	MAX_NUM_OF_EVENTS = 10000

	MIN_CATALOGUE_PRICE = 1
	MAX_CATALOGUE_PRICE = 100000

	MIN_PICK_UP_PRICE = 1
	MAX_PICK_UP_PRICE = 1000

	MIN_PRICE_PER_KM_DRIVEN = 1
	MAX_PRICE_PER_KM_DRIVEN = 100

	MIN_LENGTH_NAME = 1
	MAX_LENGTH_NAME = 40

	MIN_EVENT_TIME = 0
	MAX_EVENT_TIME = 100000

	MIN_DISTANCE_DRIVEN = 0
	MAX_DISTANCE_DRIVEN = 1000

	MIN_ACCIDENT_SEVERITY = 0
	MAX_ACCIDENT_SEVERITY = 100

	INCONSISTENT = "INCONSISTENT"

	// LOWER_CASES = "abcdefjhigklmnopqrstuvwxyz"
)

type EventType byte

const (
	Pick     EventType = 'p'
	Return   EventType = 'r'
	Accident EventType = 'a'
)

func validateName(name string) error {
	for _, item := range name {
		if unicode.IsUpper(item) || item < 'a' || item > 'z' {
			return errors.New("name must contains lower case letters only")
		}
	}
	return validateMinMax("name", len(name), MIN_LENGTH_NAME, MAX_LENGTH_NAME)
}

func evaluateError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func validateMinMax(name string, value, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be >= %d and <= %d", name, min, max)
	}
	return nil
}

type RentalEventIface interface {
	Price(carName string) int
	CheckEventType(evtType EventType) bool
	Customer() string
}

var (
	_ RentalEventIface = (*PickUpEvent)(nil)
	_ RentalEventIface = (*ReturnEvent)(nil)
	_ RentalEventIface = (*AccidentEvent)(nil)
)

type RentalEvents []RentalEventIface

func (rs RentalEvents) countByEventType(evtType EventType) int {
	count := 0
	for _, evt := range rs {
		if evt.CheckEventType(evtType) {
			count++
		}
	}
	return count
}

func (rs RentalEvents) Calculate() interface{} {
	pickCount := rs.countByEventType(Pick)

	if len(rs) <= 1 ||
		pickCount == 0 ||
		!rs[0].CheckEventType(Pick) ||
		!rs[len(rs)-1].CheckEventType(Return) {
		return INCONSISTENT
	}

	price := 0
	if pickCount == 1 {
		if rs.countByEventType(Return) != 1 {
			return INCONSISTENT
		}
		carName := rs[0].(*PickUpEvent).GetCarName()
		for _, evt := range rs {
			price += evt.Price(carName)
		}
		return price
	}

	lastPickIndex := 0

	for idx, evt := range rs {
		var procedurePrice interface{} = nil

		if evt.CheckEventType(Pick) {
			procedurePrice = rs[lastPickIndex:idx].Calculate()
			lastPickIndex = idx
		} else if idx == len(rs)-1 {
			procedurePrice = rs[lastPickIndex:].Calculate()
		}

		if procedurePrice == nil {
			continue
		}
		if procedurePrice == INCONSISTENT {
			return INCONSISTENT
		}
		price += procedurePrice.(int)
	}

	return price
}

func (s *RentalService) Audit() []string {
	res := make([]string, len(s.rentalHistory))

	idx := 0
	for customerName, rentData := range s.rentalHistory {
		res[idx] = fmt.Sprintf("%s %v", customerName, rentData.Calculate())
		idx++
	}

	sort.Strings(res)

	return res
}

func (s *RentalService) AddEvent(evt RentalEventIface) {
	customer := evt.Customer()
	s.rentalHistory[customer] = append(s.rentalHistory[customer], evt)
}

type RentalService struct {
	cars          map[string]*Car         // keys are car names. E.g: "tesla"
	rentalHistory map[string]RentalEvents // keys are unique user names. E.g "Donald_Trump"
}

type Car struct {
	Name         string
	CatalogPrice int
	PickUpPrice  int
	CostPerKm    int
}

type BaseEvent struct {
	creatAt      int
	customerName string
	eventType    EventType
	service      *RentalService // access car info with ease
}

func (b *BaseEvent) Customer() string {
	return b.customerName
}

func (b *BaseEvent) CheckEventType(evtType EventType) bool {
	return b.eventType == evtType
}

type PickUpEvent struct {
	BaseEvent
	carName string
}

func (p *PickUpEvent) GetCarName() string {
	return p.carName
}

func (p *PickUpEvent) Price(carName string) int {
	car := p.service.cars[carName]
	if car != nil {
		return car.PickUpPrice
	}
	evaluateError(fmt.Errorf("car with name=%s does not exist", carName))
	return 0
}

type ReturnEvent struct {
	DistanceDriven uint // in kilometer
	BaseEvent
}

func (p *ReturnEvent) Price(carName string) int {
	car := p.service.cars[carName]
	if car != nil {
		return int(p.DistanceDriven) * car.CostPerKm
	}
	evaluateError(fmt.Errorf("car with name=%s does not exist", carName))
	return 0
}

type AccidentEvent struct {
	Severity int // percentage
	BaseEvent
}

func (p *AccidentEvent) Price(carName string) int {
	car := p.service.cars[carName]
	if car != nil {
		floatResult := float64(p.Severity) / float64(100) * float64(car.CatalogPrice)
		return roundUp(floatResult)
	}

	evaluateError(fmt.Errorf("car with name=%s does not exist", carName))
	return 0
}

func roundUp(value float64) int {
	rounded := math.Round(value)
	if rounded < value {
		return int(rounded + 1)
	}

	return int(rounded)
}

func (s *RentalService) ListCars(numOfCars int) {
	for i := 0; i < numOfCars; i++ {
		var car Car
		_, err := fmt.Scanf("%s %d %d %d", &car.Name, &car.CatalogPrice, &car.PickUpPrice, &car.CostPerKm)
		evaluateError(err)
		evaluateError(validateName(car.Name))
		evaluateError(validateMinMax("car catalog price", car.CatalogPrice, MIN_CATALOGUE_PRICE, MAX_CATALOGUE_PRICE))
		evaluateError(validateMinMax("car pickup price", car.PickUpPrice, MIN_PICK_UP_PRICE, MAX_PICK_UP_PRICE))
		evaluateError(validateMinMax("car price per km driven", car.CostPerKm, MIN_PRICE_PER_KM_DRIVEN, MAX_PRICE_PER_KM_DRIVEN))

		s.cars[car.Name] = &car
	}
}

func (s *RentalService) LogEvents(numOfEvents int) {
	for i := 0; i < numOfEvents; i++ {
		var (
			time                    int
			customerName, eventData string
			evtType                 byte
		)

		_, err := fmt.Scanf("%d %s %c %s", &time, &customerName, &evtType, &eventData)
		evaluateError(err)
		evaluateError(validateMinMax("event time", time, MIN_EVENT_TIME, MAX_EVENT_TIME))
		evaluateError(validateName(customerName))
		if len(eventData) == 0 {
			evaluateError(fmt.Errorf("event data invalid: %s", err.Error()))
		}

		var event RentalEventIface

		switch evtType {
		case 'p':
			if s.cars[eventData] == nil {
				evaluateError(fmt.Errorf("car with name = %q does not exist in our system", eventData))
			}
			event = &PickUpEvent{BaseEvent{time, customerName, Pick, s}, eventData}

		case 'r':
			// parse distance driven
			distance, err := strconv.Atoi(eventData)
			evaluateError(err)
			evaluateError(validateMinMax("distance driven", distance, MIN_DISTANCE_DRIVEN, MAX_DISTANCE_DRIVEN))
			event = &ReturnEvent{uint(distance), BaseEvent{time, customerName, Return, s}}

		case 'a':
			severity, err := strconv.Atoi(eventData)
			evaluateError(err)
			evaluateError(validateMinMax("severity", severity, MIN_ACCIDENT_SEVERITY, MAX_ACCIDENT_SEVERITY))
			event = &AccidentEvent{severity, BaseEvent{time, customerName, Accident, s}}

		default:
			evaluateError(fmt.Errorf("invalid event type provided: %q. expected 'a' or 'r' or 'p'", evtType))
		}

		s.rentalHistory[customerName] = append(s.rentalHistory[customerName], event)
	}
}

func (s *RentalService) RecoverCrash() {
	var numOfCars, numOfEvents int
	_, err := fmt.Scanf("%d %d", &numOfCars, &numOfEvents)
	evaluateError(err)
	evaluateError(validateMinMax("num of cars", numOfCars, MIN_NUM_OF_CARS, MAX_NUM_OF_CARS))
	evaluateError(validateMinMax("num of rental events", numOfEvents, MIN_NUM_OF_EVENTS, MAX_NUM_OF_EVENTS))

	s.ListCars(numOfCars) // must go before
	s.LogEvents(numOfEvents)
}

func printResult(data []string) {
	for _, value := range data {
		fmt.Println(value)
	}
}

func main() {
	var numOfTestCases int
	_, err := fmt.Scanf("%d", &numOfTestCases)
	evaluateError(err)
	evaluateError(validateMinMax("num of test cases", numOfTestCases, MIN_NUM_OF_TESTS, MAX_NUM_OF_TESTS))

	services := make([]*RentalService, numOfTestCases)

	for i := 0; i < numOfTestCases; i++ {
		service := &RentalService{
			cars:          map[string]*Car{},
			rentalHistory: map[string]RentalEvents{},
		}
		service.RecoverCrash()
		services[i] = service
	}

	for _, s := range services {
		printResult(s.Audit())
	}
}

