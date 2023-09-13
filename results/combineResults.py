from statistics import mean

incomeFairness = {}
exp = "undefined"


with open("income.txt","r") as income:
    for line in income:
        if line[:2] == " O":
            exp = line.strip().split("-")[0]
        if line[:15] == "Income fairness":
            values = incomeFairness.get(exp,[])
            values.append(float(line.split(":")[1]))
            incomeFairness[exp] = values


for (exp, values) in incomeFairness.items():
    print(exp,mean(values), max(values))

print("WORK:")

workFairness = {}
forwardworkFairness = {}
exp = "undefined"

with open("work.txt","r") as work:
    for line in work:
        if line[:2] == " O":
            exp = line.strip().split("-")[0]
        if line[:12] == "Workfairness":
            values = workFairness.get(exp,[])
            values.append(float(line.split(":")[1]))
            workFairness[exp] = values
        if line[:19] == "Forwardworkfairness":
            values = workFairness.get(exp+"Fwd",[])
            values.append(float(line.split(":")[1]))
            workFairness[exp+"Fwd"] = values
        if line[:19] == "Storageworkfairness":
            values = workFairness.get(exp+"Store",[])
            values.append(float(line.split(":")[1]))
            workFairness[exp+"Store"] = values
        if line[:3] == "Max":
            values = workFairness.get(exp+"Max",[])
            values.append(float(line.split(":")[1].split(",")[0]))
            workFairness[exp+"Max"] = values
            values = workFairness.get(exp+"Max!O",[])
            values.append(float(line.split(":")[1].split(",")[1]))
            workFairness[exp+"Max!O"] = values
            values = workFairness.get(exp+"Med",[])
            values.append(float(line.split(":")[1].split(",")[2]))
            workFairness[exp+"Med"] = values


for (exp, values) in workFairness.items():
    if exp[-5:] == "Store":
        print(exp,mean(values), max(values))
for (exp, values) in workFairness.items():
    if exp[-3:] == "Fwd":
        print(exp,mean(values), max(values))
for (exp, values) in workFairness.items():
    if exp[-3:] != "ore" and exp[-3:] != "Fwd":
        print(exp,mean(values), max(values))

print("HOPS:")
hopsincome = {}
exp = "undefined"
avgfwdreward = {}
stdfwdreward = {}

with open("hops.txt","r") as hops:
    for line in hops:
        if line[:2] == " O":
            exp = line.strip().split("-")[0]
        if line[:4] == "Hop:":
            hop = line.strip().split(" ")[1]
            income = line.strip().split(" ")[-1]
            hopdict = hopsincome.get(exp,{})
            values = hopdict.get(hop, [])
            values.append(float(income))
            hopdict[hop]=values
            hopsincome[exp]=hopdict
        if line[:12] == "Mean and std":
            vals = line.strip().split(":")[1].split(",")
            avg = vals[0]
            std = vals[1]
            avgs = avgfwdreward.get(exp,[])
            avgs.append(float(avg))
            avgfwdreward[exp] = avgs
            stds = stdfwdreward.get(exp, [])
            stds.append(float(std))
            stdfwdreward[exp] = stds




for (exp, hopdict) in hopsincome.items():
    print(exp, end=", ")
    incomelist = [-1,-1,-1,-1,-1,-1,-1,-1,-1]
    total = 0
    for (hop, values) in hopdict.items():
        meanv = mean(values)
        
        if hop == "storer":
            incomelist[-1]= meanv
        else:
            # OUPS does not include storer value in total
            total += meanv
            incomelist[int(hop)] = meanv

    for val in incomelist:
        print("{:4f}".format(val/total), end=", ")
    print()

print("mean forwads reward")
for exp in avgfwdreward:
    print("{}, {:5f}, {:5f}".format(exp, mean(avgfwdreward[exp]), mean(stdfwdreward[exp])))

print("PayHOPS:")
hopsincome = {}
exp = "undefined"
avgfwdreward = {}
stdfwdreward = {}

with open("hopPays.txt","r") as hops:
    for line in hops:
        if line[:2] == " O":
            exp = line.strip().split("-")[0]
        if line[:4] == "Hop:":
            hop = line.strip().split(" ")[1]
            income = line.strip().split(" ")[-1]
            hopdict = hopsincome.get(exp,{})
            values = hopdict.get(hop, [])
            values.append(float(income))
            hopdict[hop]=values
            hopsincome[exp]=hopdict
        if line[:12] == "Mean and std":
            vals = line.strip().split(":")[1].split(",")
            avg = vals[0]
            std = vals[1]
            avgs = avgfwdreward.get(exp,[])
            avgs.append(float(avg))
            avgfwdreward[exp] = avgs
            stds = stdfwdreward.get(exp, [])
            stds.append(float(std))
            stdfwdreward[exp] = stds




for (exp, hopdict) in hopsincome.items():
    print(exp, end=", ")
    incomelist = [-1,-1,-1,-1,-1,-1,-1,-1,-1]
    total = 0
    for (hop, values) in hopdict.items():
        meanv = mean(values)
        
        if hop == "storer":
            incomelist[-1]= meanv
        else:
            # OUPS does not include storer value in total
            total += meanv
            incomelist[int(hop)] = meanv

    for val in incomelist:
        print("{:4f}".format(val/total), end=", ")
    print()

print("mean forwar reward")
for exp in avgfwdreward:
    print("{}, {:5f}, {:5f}".format(exp, mean(avgfwdreward[exp]), mean(stdfwdreward[exp])))
