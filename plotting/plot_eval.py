import json
import matplotlib.pyplot as plt
import matplotlib

cmap = matplotlib.cm.get_cmap('Pastel1')
algos = [None] * 5
times = [None] * 5
times_faster = [None] * 5
pqpops = [None] * 5
pqpops_percent = [None] * 5
colors = ['yellowgreen', 'deepskyblue', 'darkorchid', 'orange', 'firebrick']
data = {}

with open('../data/evaluation/wf_360_360_2020-10-02_12-08-17.json') as file:
    data = json.load(file)

for i, x in enumerate(data['Results']):
    algos[i] = x
    times[i] = round(float(data['Results'][x]['Time']['AVG'].replace("ms", "")), 3)
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
ax.bar(algos1, times, width=0.5, color=colors1)
ax.set_title('Uniform Grid (n = 129.600)')

for i in range(len(times)):
    if i != 0:
        plt.annotate(str(times[i]) + "\n(x" + str(times_faster[i]) + ")",
                     xy=(algos1[i], times[i]), ha='center', va='bottom')
    else:
        plt.annotate(str(times[i]), xy=(algos1[i], times[i]), ha='center', va='bottom')
plt.savefig("Uniform Grid (n = 129.600) speed.jpg")
plt.show()

#
# PQ-Pops
#
pqpops, algos2, colors2, pqpops_percent = zip(
    *sorted(zip(pqpops, algos, colors, pqpops_percent), reverse=True))

fig, ax = plt.subplots()
ax.set_ylabel('Average priority queue pops')
ax.bar(algos2, pqpops, width=0.5, color=colors2)
ax.set_title('Uniform Grid (n = 129.600)')

for i in range(len(pqpops)):
    if i != 0:
        plt.annotate(str(pqpops[i]) + "\n(" + str(pqpops_percent[i]) + "%)",
                     xy=(algos2[i], pqpops[i]), ha='center', va='bottom')
    else:
        plt.annotate(str(pqpops[i]), xy=(algos2[i], pqpops[i]), ha='center', va='bottom')
plt.savefig("Uniform Grid (n = 129.600) pqpops.jpg")
plt.show()
