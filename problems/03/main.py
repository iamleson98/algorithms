def parse_lines(lines: list):
  tests = []
  car_histories = {}
  car_prices = {}
  spy_histories = {}
  test_cases = 0
  event_num = 0
  event_count = 0
  for l in lines:
    l_split = l.strip().split()
    if len(l_split) == 1:
      test_cases = int(l.strip())
    elif len(l_split) == 2:
      event_num = int(l_split[1])
      event_count = 0
    elif not l_split[0].isdigit():
      car_prices[l_split[0]] = {'catalog': int(l_split[1]), 'pickup': int(l_split[2]), 'mileague': int(l_split[3])}
    else:
      action = l_split[2]
      spy = l_split[1]
      time = int(l_split[0])
      if action == 'p':
        car_histories[l_split[2]] = car_histories.get(l_split[2], {})
        car_histories[l_split[2]][time] = {'action': action, 'spy': spy}
        spy_histories[spy] = spy_histories.get(spy, {})
        spy_histories[spy][time] = {'action': action, 'car': l_split[2]}
      elif action == 'a':
        spy_histories[spy] = spy_histories.get(spy, {})
        spy_histories[spy][time] = {'action': action, 'damage': int(l_split[2])}
      else:
        spy_histories[spy] = spy_histories.get(spy, {})
        spy_histories[spy][time] = {'action': action, 'mileague': int(l_split[2])}
      event_count += 1
      if event_count == event_num:
        event_count = 0
        tests.append({'car_histories': car_histories, 'car_prices': car_prices, 'spy_histories': spy_histories})
        car_histories, car_prices, spy_histories = {}, {}, {}
  return tests
      

def solve(test_case: dict):
  spy_histories = test_case['spy_histories']
  log_stack = {s: [] for s in spy_histories}
  verdicts = {s: 0 for s in spy_histories}
  car_prices = test_case['car_prices']
  for spy, timeline in spy_histories.items():
    if str(verdicts.get(spy, '')) == 'INCONSISTENT':
      continue
    for time in sorted(timeline.keys()):
      if timeline[time]['action'] == 'p':
        # If already has a car
        if log_stack[spy]:
          verdicts[spy] = 'INCONSISTENT'
          break
        # If not, add the car to stack
        log_stack[spy].append(timeline[time]['car'])
        verdicts[spy] += car_prices[car]['pickup']
      elif timeline[time]['action'] in ['a', 'r']:
        # If has no car
        if not log_stack[spy]:
          verdicts[spy] = 'INCONSISTENT'
          break
        car = log_stack[spy][-1]
        if timeline[time]['action']= 'r':
          # Remove car from stack
          log_stack[spy].pop()
          verdicts[spy] += car_prices[car]['mileague'] * timeline[time]['mileague']
        else:
          repair_cost = car_prices[car]['catalog'] * timeline[time]['damage'] / 100
          # Always round up in case of float
          if int(repair_cost) != repair_cost:
            repair_cost = int(repair_cost) + 1
          verdicts[spy] += repair_cost
  # Check if anyone still has car
  for spy, stack in log_stack.items():
    if stack:
      verdicts[spy] = 'INCONSISTENT'
  return verdicts
          
          
