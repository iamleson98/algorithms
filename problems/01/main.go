package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

const (
	MAX_NUMBER_OF_DISHES = 1000
	MIN_NUMBER_OF_DISHES = 1

	MAX_NUMBER_OF_INGREDIENTS = 20
	MIN_NUMBER_OF_INGREDIENTS = 1

	MIN_NUMBER_OF_STANDARD_PORTIONS = 1
	MAX_NUMBER_OF_STANDARD_PORTIONS = 12

	MIN_NUMBER_OF_DESIRED_PORTIONS = 1
	MAX_NUMBER_OF_DESIRED_PORTIONS = 1000

	MAX_LENGTH_INGREDIENT_NAME = 20
	MIN_LENGTH_INGREDIENT_NAME = 1
)

type Dish struct {
	Ingredients []Ingredient
	// calculate desired weight amount of each ingredient for output
	Formular func(ingre Ingredient) float64
}

type Ingredient struct {
	Name          string
	WeightGr      float64
	WeightPercent float64
}

// evaluateError checks given err. If non-nil log to the standard output and exit the program
func evaluateError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// validateMinMax validates if given value is within [min, max]
func validateMinMax(value, min, max int) error {
	if value < min || value > max {
		var format string = "error: expected value V to be: %d <= V <= %d"

		kind := reflect.ValueOf(value).Kind()
		if kind == reflect.Float32 || kind == reflect.Float64 {
			format = "error: expected value V to be: %f <= V <= %f"
		}
		return fmt.Errorf(format, min, max)
	}

	return nil
}

func parseDishes(numOfDishes int) []Dish {
	var dishes = make([]Dish, numOfDishes)

	for i := 0; i < numOfDishes; i++ {
		var numOfIngredients, numOfStandardPortions, numOfDesiredPortions int
		dashHasMainIngredient := false

		_, err := fmt.Scanf("%d %d %d", &numOfIngredients, &numOfStandardPortions, &numOfDesiredPortions)
		evaluateError(err)

		evaluateError(validateMinMax(numOfIngredients, MIN_NUMBER_OF_INGREDIENTS, MAX_NUMBER_OF_INGREDIENTS))
		evaluateError(validateMinMax(numOfStandardPortions, MIN_NUMBER_OF_STANDARD_PORTIONS, MAX_NUMBER_OF_STANDARD_PORTIONS))
		evaluateError(validateMinMax(numOfDesiredPortions, MIN_NUMBER_OF_DESIRED_PORTIONS, MAX_NUMBER_OF_DESIRED_PORTIONS))

		dishes[i].Ingredients = make([]Ingredient, numOfIngredients)

		for j := 0; j < numOfIngredients; j++ {
			var ingredient Ingredient

			_, err = fmt.Scanf("%s %f %f", &ingredient.Name, &ingredient.WeightGr, &ingredient.WeightPercent)
			evaluateError(err)

			err = validateMinMax(len(ingredient.Name), MIN_LENGTH_INGREDIENT_NAME, MAX_LENGTH_INGREDIENT_NAME)
			if err != nil || strings.Contains(ingredient.Name, " ") {
				evaluateError(fmt.Errorf("error: ingredient name must has length >= %d, <= %d and contains no space character", MIN_LENGTH_INGREDIENT_NAME, MAX_LENGTH_INGREDIENT_NAME))
			}

			if fmt.Sprintf("%.1f", ingredient.WeightPercent) == "100.0" {
				dashHasMainIngredient = true
				dishes[i].Formular = func(ingre Ingredient) float64 {
					return ingre.WeightPercent * (float64(numOfDesiredPortions) / float64(numOfStandardPortions) * ingredient.WeightGr) / 100
				}
			}

			dishes[i].Ingredients[j] = ingredient
		}

		if !dashHasMainIngredient {
			evaluateError(errors.New("this dish has no main ingredient"))
		}
	}

	return dishes
}

func outPut(dishes []Dish) {
	for i, dish := range dishes {
		fmt.Printf("Recipe # %d\n", i+1)
		for _, ingre := range dish.Ingredients {
			fmt.Printf("%s %.1f\n", ingre.Name, dish.Formular(ingre))
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}

func main() {
	var numberOfDishes int
	_, err := fmt.Scanf("%d", &numberOfDishes)
	evaluateError(err)
	validateMinMax(numberOfDishes, MIN_NUMBER_OF_DISHES, MAX_NUMBER_OF_DISHES)

	dishes := parseDishes(numberOfDishes)
	outPut(dishes)
}
