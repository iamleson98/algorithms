import math

def parse_lines(lines: list):
  total_cost = {}
  car_prices = {}
  spies = {}
  inconsistent_flag = {}
  event_dict = {'p':0, 'a':1, 'r':2}
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
      if l_split[1] not in spies and l_split[1] not in inconsistent_flag:
        if l_split[2] == 'p':
          spies[l_split[1]] = {'car': [], 'time': [], 'events':[], 'costs':[]}
        else:
          inconsistent_flag[l_split[1]] = 'True'
      # else:
        # if l_split[2] == 'p':
        # inconsistent_flag[l_split[1]] = 'inconsistent'
      if l_split[1] in inconsistent_flag:
        total_cost[l_split[1]] = 'Inconsistent'
      else:
      
        spies[l_split[1]]['time'].append(l_split[0])
        spies[l_split[1]]['events'].append(event_dict[l_split[2]])
        
        if l_split[2] == 'p':
          car_type = l_split[-1]
          spies[l_split[1]]['car'].append(l_split[-1])
          spies[l_split[1]]['costs'].append(car_prices[car_type]['pickup'])
        elif l_split[2] == 'r':
          car_type = spies[l_split[1]]['car'][-1]
          running_cost = car_prices[car_type]['mileague']*l_split[3]
          spies[l_split[1]]['costs'].append(running_cost)
        elif l_split[3] == 'a':
          car_type = spies[l_split[1]]['car'][-1]
          accident_cost = math.ceil(car_prices[car_type]['catalog']*l_split[3]/100)
          spies[l_split[1]]['costs'].append(accident_cost)
        
  #done reading
  for spies_name in spies:
    # if spies_name not in inconsistent_flag:
    spies_events = spies[spies_name]['events']
    unique_element = set(spies[spies_events])
    count_unique = [spies_events.count(elm) for elm in unique_element if spies_events.count(elm)>1 and elm != 1]
    if len(count_unique)>0:
      total_cost[spies_name] = 'Inconsistent'
    else:
      if spies_events[0] == 0 and spies_events[-1] == 2: 
        total_cost[spies_name] = sum(spies[spies_name]['costs'])
      else:
          total_cost[spies_name] = 'Inconsistent'
  
  total_cost = dict(sorted(total_cost.items()))
  return total_cost
                  
          
          
