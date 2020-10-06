import json
import matplotlib.pyplot as plt
import matplotlib
import argparse

parser = argparse.ArgumentParser()
parser.add_argument('filename', metavar='f', type=str, nargs='+',
                    help='filename')
args = parser.parse_args()
filename = args.filename[0]

size = 0
try:
    x = filename.split("_")[1]
    y = filename.split("_")[2]
    size = int(x) * int(y)
except Exception:
    print("Filename not supported.")
    exit()

cmap = matplotlib.cm.get_cmap('Pastel1')
algos = [None] * 5
times = [None] * 5
times_faster = [None] * 5
pqpops = [None] * 5
pqpops_percent = [None] * 5
colors = ['yellowgreen', 'deepskyblue', 'darkorchid', 'orange', 'firebrick']
data = {}
# ../data/evaluation/wf_360_360_2020-10-02_12-08-17.json

try:
    with open(filename) as file:
        data = json.load(file)
        # print(json.dumps(data, indent="  "))
except Exception:
    print("File not found.")
    exit()


for i, x in enumerate(data['Results']):
    algos[i] = x
    timestring = data['Results'][x]['Time']['AVG']
    if timestring.endswith("ms"):
        times[i] = round(float(data['Results'][x]['Time']['AVG'].replace('ms', '')), 3)
    else:
        timestring = timestring.replace("s", "")
        timestring = round(float(timestring) * 1000)
        times[i] = timestring
    times_faster[i] = round(float(data['Results'][x]['Time']['TimesFaster']), 3)
    pqpops[i] = data['Results'][x]['PQPops']['AVG']
    pqpops_percent[i] = data['Results'][x]['PQPops']['Percent']

#
# Times
#
times, algos1, colors1, times_faster = zip(
    *sorted(zip(times, algos, colors, times_faster), reverse=True))

fig, ax = plt.subplots()
ax.set_ylabel('Average time in ms')
ax.set_xlabel('Averaging over %s runs with %s CPU cores.'
              % (data["Parameters"]["Run Count"], data["Parameters"]["CPU Cores"]))
ax.bar(algos1, times, width=0.5, color=colors1)
ax.set_title('Uniform Grid (n = %d)' % size)

for i in range(len(times)):
    if i != 0:
        plt.annotate(str(times[i]) + '\n(x' + str(times_faster[i]) + ')',
                     xy=(algos1[i], times[i]), ha='center', va='bottom')
    else:
        plt.annotate(str(times[i]), xy=(algos1[i], times[i]), ha='center', va='bottom')
plt.savefig('Uniform Grid (n = %d) speed.jpg' % size)
plt.show()

#
# PQ-Pops
#
pqpops, algos2, colors2, pqpops_percent = zip(
    *sorted(zip(pqpops, algos, colors, pqpops_percent), reverse=True))

fig, ax = plt.subplots()
ax.set_ylabel('Average priority queue pops')
ax.set_xlabel('Averaging over %s runs with %s CPU cores.'
              % (data["Parameters"]["Run Count"], data["Parameters"]["CPU Cores"]))
ax.bar(algos2, pqpops, width=0.5, color=colors2)
ax.set_title('Uniform Grid (n = %d)' % size)

for i in range(len(pqpops)):
    if i != 0:
        plt.annotate(str(pqpops[i]) + '\n(' + str(pqpops_percent[i]) + '%)',
                     xy=(algos2[i], pqpops[i]), ha='center', va='bottom')
    else:
        plt.annotate(str(pqpops[i]), xy=(algos2[i], pqpops[i]), ha='center', va='bottom')
plt.savefig('Uniform Grid (n = %d) pqpops.jpg' % size)
plt.show()
