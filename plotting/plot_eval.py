import json
import matplotlib
import matplotlib.pyplot as plt
import argparse
# plt.style.use('dark_background')

parser = argparse.ArgumentParser()
parser.add_argument('filename', type=str, nargs='+',
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

if "bg" in filename:
    title = 'Basic Grid (n = %d)' % size
else:
    title = 'Uniform Grid (n = %d)' % size

cmap = matplotlib.cm.get_cmap('Pastel1')
algos = [None] * 5
equality = [None] * 5
times = [None] * 5
times_faster = [None] * 5
pqpops = [None] * 5
pqpops_percent = [None] * 5
colors = ['yellowgreen', 'deepskyblue', 'darkorchid', 'orange', 'firebrick']
data = {}

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
    equality[i] = [data['Results'][x]['Length']['Shorter'], data['Results']
                   [x]['Length']['Equal'], data['Results'][x]['Length']['Longer']]

#
# Equality
#
fig, axs = plt.subplots(2, 2, figsize=(12, 8))

axs[0, 0].set_title(algos[3])
axs[0, 0].pie(equality[3], colors=colors, shadow=True, startangle=0, pctdistance=1.3,
              autopct=lambda p: '{:.1f}%'.format(round(p)) if p > 0 else '')

axs[0, 1].set_title(algos[2])
axs[0, 1].pie(equality[2], colors=colors, shadow=True, startangle=0, pctdistance=1.3,
              autopct=lambda p: '{:.1f}%'.format(round(p)) if p > 0 else '')

axs[1, 0].set_title(algos[0])
axs[1, 0].pie(equality[0], colors=colors, shadow=True, startangle=0, pctdistance=1.3,
              autopct=lambda p: '{:.1f}%'.format(round(p)) if p > 0 else '')

axs[1, 1].set_title(algos[1])
axs[1, 1].pie(equality[1], colors=colors, shadow=True, startangle=-90, pctdistance=1.3,
              autopct=lambda p: '{:.1f}%'.format(round(p)) if p > 0 else '')

fig.legend(['Shorter', 'Equal', 'Longer'], loc="center")
fig.suptitle(title)
plt.tight_layout(pad=0.4, w_pad=0.5, h_pad=1.0)
plt.savefig('%s equality.jpg' % title)
plt.show()

#
# Times
#
times, algos1, colors1, times_faster = zip(
    *sorted(zip(times, algos, colors, times_faster), reverse=True))

fig, ax = plt.subplots(figsize=(12, 8))
ax.set_ylabel('Average time in ms')
ax.set_xlabel('Averaging over %s runs with %s CPU cores.'
              % (data["Parameters"]["Run Count"], data["Parameters"]["CPU Cores"]))
ax.bar(algos1, times, color=colors1)
# ax.bar(algos1, times, width=0.5, color=colors1)
ax.set_title('Speed for %s' % title)

for i in range(len(times)):
    if i != 0:
        plt.annotate(str(times[i]) + '\n(x' + str(times_faster[i]) + ')',
                     xy=(algos1[i], times[i]), ha='center', va='bottom')
    else:
        plt.annotate(str(times[i]), xy=(algos1[i], times[i]), ha='center', va='bottom')
plt.tight_layout()
plt.savefig('%s speed.jpg' % title)
plt.show()

#
# PQ-Pops
#
pqpops, algos2, colors2, pqpops_percent = zip(
    *sorted(zip(pqpops, algos, colors, pqpops_percent), reverse=True))

fig, ax = plt.subplots(figsize=(12, 8))
ax.set_ylabel('Average priority queue pops')
ax.set_xlabel('Averaging over %s runs with %s CPU cores.'
              % (data["Parameters"]["Run Count"], data["Parameters"]["CPU Cores"]))
ax.bar(algos2, pqpops, color=colors2)
# ax.bar(algos2, pqpops, width=0.5, color=colors2)
ax.set_title('PQ-Pops for %s' % title)

for i in range(len(pqpops)):
    if i != 0:
        plt.annotate(str(pqpops[i]) + '\n(' + str(pqpops_percent[i]) + '%)',
                     xy=(algos2[i], pqpops[i]), ha='center', va='bottom')
    else:
        plt.annotate(str(pqpops[i]), xy=(algos2[i], pqpops[i]), ha='center', va='bottom')
plt.tight_layout()
plt.savefig('%s pqpops.jpg' % title)
plt.show()
