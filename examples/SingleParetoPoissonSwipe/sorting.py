import csv
import numpy as np
import sys
import os
import subprocess




def bu(r_dict):
    return (r_dict["LINK-AvgC"]/r_dict["LINK-CBuff"] - r_dict["GEWI-AvgC"]/r_dict["GEWI-CBuff"] )

def mwt(r_dict):
    return r_dict["LINK-MWT"]-r_dict["GEWI-MWT"]

def drop(r_dict):
    return (r_dict["LINK-Droppd"]/r_dict["LINK-Recv"]-r_dict["GEWI-Droppd"]/r_dict["GEWI-Recv"])

def enbuff(r_dict):
    return (r_dict["GEWI-AvgE"]/r_dict["GEWI-EnBuff"])

def enbuff(r_dict):
    return (r_dict["GEWI-Transm"]/r_dict["GEWI-EnBuff"])

def analyze(bu, in_file ,out_file):
    ResultsDict = {}
    with open('results/'+ in_file + ".csv") as file:
        csv_reader = csv.reader(file, delimiter=',')
        names = []
        line_count = 0
        for row in csv_reader:
            if line_count == 0:
                names = row
                line_count += 1
                continue
            r_dict = {}
            for i,r in enumerate(row):
                r_dict[names[i]]=float(r)
            line_count += 1

            if r_dict["PARETO-NodeCount"] not in ResultsDict:
                ResultsDict[r_dict["PARETO-NodeCount"]] = {
                    "BufferUtil_1_AVG_M=2":None,
                    "BufferUtil_1_AVG_CALC_M=2":[],
                    "BufferUtil_2_AVG_M=2":None,
                    "BufferUtil_2_AVG_CALC_M=2":[],
                    "BufferUtil_3_AVG_M=2":None,
                    "BufferUtil_3_AVG_CALC_M=2":[],
                    "BufferUtil_1_STD_M=2":None,
                    "BufferUtil_2_STD_M=2":None,
                    "BufferUtil_3_STD_M=2":None,

                    "BufferUtil_1_AVG_M=16":None,
                    "BufferUtil_1_AVG_CALC_M=16":[],
                    "BufferUtil_2_AVG_M=16":None,
                    "BufferUtil_2_AVG_CALC_M=16":[],
                    "BufferUtil_3_AVG_M=16":None,
                    "BufferUtil_3_AVG_CALC_M=16":[],
                    "BufferUtil_1_STD_M=16":None,
                    "BufferUtil_2_STD_M=16":None,
                    "BufferUtil_3_STD_M=16":None,

                    "BufferUtil_1_AVG_M=128":None,
                    "BufferUtil_1_AVG_CALC_M=128":[],
                    "BufferUtil_2_AVG_M=128":None,
                    "BufferUtil_2_AVG_CALC_M=128":[],
                    "BufferUtil_3_AVG_M=128":None,
                    "BufferUtil_3_AVG_CALC_M=128":[],
                    "BufferUtil_1_STD_M=128":None,
                    "BufferUtil_2_STD_M=128":None,
                    "BufferUtil_3_STD_M=128":None,
                }
            if r_dict["GEWI-Mult"] == 2:
                if r_dict["PARETO-Shape_on"] < 1.33:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_1_AVG_CALC_M=2"].append(bu(r_dict))
                elif r_dict["PARETO-Shape_on"] >= 1.33 and r_dict["PARETO-Shape_on"] <=1.66:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_2_AVG_CALC_M=2"].append(bu(r_dict))
                elif r_dict["PARETO-Shape_on"] > 1.66:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_3_AVG_CALC_M=2"].append(bu(r_dict))


            if r_dict["GEWI-Mult"] == 16:
                if r_dict["PARETO-Shape_on"] < 1.33:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_1_AVG_CALC_M=16"].append(bu(r_dict))
                elif r_dict["PARETO-Shape_on"] >= 1.33 and r_dict["PARETO-Shape_on"] <=1.66:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_2_AVG_CALC_M=16"].append(bu(r_dict))
                elif r_dict["PARETO-Shape_on"] > 1.66:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_3_AVG_CALC_M=16"].append(bu(r_dict))

            if r_dict["GEWI-Mult"] == 128:

                if r_dict["PARETO-Shape_on"] < 1.33:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_1_AVG_CALC_M=128"].append(bu(r_dict))
                elif r_dict["PARETO-Shape_on"] >= 1.33 and r_dict["PARETO-Shape_on"] <=1.66:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_2_AVG_CALC_M=128"].append(bu(r_dict))
                elif r_dict["PARETO-Shape_on"] > 1.66:
                    ResultsDict[r_dict["PARETO-NodeCount"]]["BufferUtil_3_AVG_CALC_M=128"].append(bu(r_dict))

    for n in ResultsDict.keys():
        ResultsDict[n]["BufferUtil_1_AVG_M=2"] = np.average(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=2"])
        ResultsDict[n]["BufferUtil_2_AVG_M=2"] = np.average(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=2"])
        ResultsDict[n]["BufferUtil_3_AVG_M=2"] = np.average(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=2"])

        ResultsDict[n]["BufferUtil_1_STD_M=2"] = np.std(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=2"])
        ResultsDict[n]["BufferUtil_2_STD_M=2"] = np.std(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=2"])
        ResultsDict[n]["BufferUtil_3_STD_M=2"] = np.std(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=2"])

        ResultsDict[n]["BufferUtil_1_AVG_M=16"] = np.average(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=16"])
        ResultsDict[n]["BufferUtil_2_AVG_M=16"] = np.average(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=16"])
        ResultsDict[n]["BufferUtil_3_AVG_M=16"] = np.average(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=16"])
        ResultsDict[n]["BufferUtil_1_STD_M=16"] = np.std(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=16"])
        ResultsDict[n]["BufferUtil_2_STD_M=16"] = np.std(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=16"])
        ResultsDict[n]["BufferUtil_3_STD_M=16"] = np.std(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=16"])

        ResultsDict[n]["BufferUtil_1_AVG_M=128"] = np.average(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=128"])
        ResultsDict[n]["BufferUtil_2_AVG_M=128"] = np.average(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=128"])
        ResultsDict[n]["BufferUtil_3_AVG_M=128"] = np.average(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=128"])

        ResultsDict[n]["BufferUtil_1_STD_M=128"] = np.std(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=128"])
        ResultsDict[n]["BufferUtil_2_STD_M=128"] = np.std(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=128"])
        ResultsDict[n]["BufferUtil_3_STD_M=128"] = np.std(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=128"])

        del(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=2"])
        del(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=2"])
        del(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=2"])

        del(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=16"])
        del(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=16"])
        del(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=16"])

        del(ResultsDict[n]["BufferUtil_1_AVG_CALC_M=128"])
        del(ResultsDict[n]["BufferUtil_2_AVG_CALC_M=128"])
        del(ResultsDict[n]["BufferUtil_3_AVG_CALC_M=128"])

    out = out_file + ".csv"
    if os.path.exists(out):
        os.remove(out)
    with open(out, "a") as f:
        writer = csv.writer(f)
        keys = list(ResultsDict.keys())
        keys.sort()
        writer.writerow(["PARETO-NodeCount",
                        "BufferUtil_1_AVG_M=2", "BufferUtil_1_STD_M=2",
                        "BufferUtil_2_AVG_M=2", "BufferUtil_2_STD_M=2",
                        "BufferUtil_3_AVG_M=2", "BufferUtil_3_STD_M=2",

                        "BufferUtil_1_AVG_M=16", "BufferUtil_1_STD_M=16",
                        "BufferUtil_2_AVG_M=16", "BufferUtil_2_STD_M=16",
                        "BufferUtil_3_AVG_M=16", "BufferUtil_3_STD_M=16",

                        "BufferUtil_1_AVG_M=128", "BufferUtil_1_STD_M=128",
                        "BufferUtil_2_AVG_M=128", "BufferUtil_2_STD_M=128",
                        "BufferUtil_3_AVG_M=128", "BufferUtil_3_STD_M=128"])
        for key in keys:
            r = ResultsDict[key]
            arr = [key,
                r['BufferUtil_1_AVG_M=2'],r['BufferUtil_1_STD_M=2'],
                r['BufferUtil_2_AVG_M=2'],r['BufferUtil_2_STD_M=2'],
                r['BufferUtil_3_AVG_M=2'],r['BufferUtil_3_STD_M=2'],
                r['BufferUtil_1_AVG_M=16'],r['BufferUtil_1_STD_M=16'],
                r['BufferUtil_2_AVG_M=16'],r['BufferUtil_2_STD_M=16'],
                r['BufferUtil_3_AVG_M=16'],r['BufferUtil_3_STD_M=16'],
                r['BufferUtil_1_AVG_M=128'],r['BufferUtil_1_STD_M=128'],
                r['BufferUtil_2_AVG_M=128'],r['BufferUtil_2_STD_M=128'],
                r['BufferUtil_3_AVG_M=128'],r['BufferUtil_3_STD_M=128'],
                ]
            writer.writerow(arr)


if __name__=="__main__":
    in_file = sys.argv[1]
    print("Fetching data")
    scp_call = ["scp","root@75.119.131.45:GO-DTS/examples/ParetoSwipe/results/"+in_file+".csv","results/"+in_file+".csv"]
    subprocess.call(scp_call)
    print("Analyzing data")
    analyze(bu, in_file, ".butil")
    analyze(mwt, in_file, ".mwt")
    analyze(drop, in_file, ".drop")
    analyze(enbuff, in_file, ".enbuff")
    print("Plotting")
    subprocess.call(["gnuplot", "MULTIPLOT.plt"])
