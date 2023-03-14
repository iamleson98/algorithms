import math

testCases = input()
for testCase in testCases:
	nextLine = input().split()
	cars = int(nextLine[0])
	events = int(nextLine[1])
	carDict = {}
	for i in range(cars):
		carLine = input().split()
		carDict[carLine[0]] = (int(carLine[1]), int(carLine[2]), int(carLine[3]))
	result = {}
	for i in range(events):
		time, spy, e, eventValue = input().split()
		if spy not in result:
			result[spy] = [0, 0, '']
		if result[spy][0] < 0:
			continue
		if e == 'p':
			if result[spy][0] == 1:
				result[spy][0] = -1
			else:
				result[spy][1] += carDict[eventValue][1]
				result[spy][2] = eventValue
				result[spy][0] = 1
		elif e == 'r':
			if result[spy][0] == 0:
				result[spy] = -1
			else:
				result[spy][1] += carDict[result[spy][2]][2]*int(eventValue)
				result[spy][0] = 0
		else:
			if result[spy][0] == 0:
				result[spy][0] = -1
			else:
				result[spy][1] += math.ceil(carDict[result[spy][2]][0]*int(eventValue)/100)
	for spy in sorted(result.keys()):
		print(spy, result[spy][1] if result[spy][0] == 0 else 'INCONSISTENT')
