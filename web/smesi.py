import requests
import time
from matplotlib import pyplot as plt

TANK0 = 'Бак0'
TANK1 = 'Бак1'
TANK2 = 'Бак2'

DB1 = 'ДБ1'
DB2 = 'ДБ2'

DU0 = 'ДУ0'
DU1 = 'ДУ1'
DU2 = 'ДУ2'

KL0 = 'Кл0'
KL1 = 'Кл1'
KL2 = 'Кл2'
KL3 = 'Кл3'

def get_sensor_info(sensor_name):
    r = requests.get("http://localhost:8080/sensor/"+sensor_name)
    return r.json()

def get_tank_info(tank_name):
    r = requests.get("http://localhost:8080/tank/"+tank_name)
    return r.json()

def open_valve(valve_name):
    print('Open valve: ' + valve_name)
    r = requests.post("http://localhost:8080/valve/open/"+valve_name)

def close_valve(valve_name):
    print('Close valve: ' + valve_name)
    r = requests.post("http://localhost:8080/valve/close/"+valve_name)

def toggle_rotor():
    print('Toggle rotor')
    r = requests.post("http://localhost:8080/drp")

time_info = []
tank0_info = []
tank1_info = []
tank2_info = []
#sensors_info = []

def get_info(start_time):
    time_info.append(time.time()-start_time)

    t0 = get_tank_info(TANK0)
    t1 = get_tank_info(TANK1)
    t2 = get_tank_info(TANK2)
    tank0_info.append(t0['CurVolume']/t0['MaxVolume'])
    tank1_info.append(t1['CurVolume']/t1['MaxVolume'])
    tank2_info.append(t2['CurVolume']/t2['MaxVolume'])

def work_cycle(start_time):
    get_info(start_time)

    open_valve(KL1)
    while not get_sensor_info(DU1)['Value']:
        get_info(start_time)
        time.sleep(0.01)
    close_valve(KL1)

    get_info(start_time)

    open_valve(KL2)
    while not get_sensor_info(DU2)['Value']:
        get_info(start_time)
        time.sleep(0.01)
    close_valve(KL2)

    get_info(start_time)

    toggle_rotor()
    time.sleep(2)
    toggle_rotor()

    get_info(start_time)

    open_valve(KL0)
    while get_sensor_info(DU0)['Value']:
        get_info(start_time)
        time.sleep(0.01)
    close_valve(KL0)

    get_info(start_time)


def main():
    start_time = time.time()
    
    for _ in range(10):
        work_cycle(start_time)
    
    plt.plot(time_info, tank0_info)
    plt.ylabel('Объем вещества в баке 0')
    plt.xlabel('Время, c')
    plt.grid(True)
    plt.title('Бак 0')
    plt.show()

    plt.plot(time_info, tank1_info)
    plt.ylabel('Объем вещества в баке 1')
    plt.xlabel('Время, c')
    plt.grid(True)
    plt.title('Бак 1')
    plt.show()

    plt.plot(time_info, tank2_info)
    plt.ylabel('Объем вещества в баке 2')
    plt.xlabel('Время, c')
    plt.grid(True)
    plt.title('Бак 2')
    plt.show()

main()