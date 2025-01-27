import pandas as pd
import matplotlib.pyplot as plt
import json
from sklearn.linear_model import LinearRegression
import numpy as np
import argparse

def main():
    rawdata = ""
    trace = ""

    parser = argparse.ArgumentParser(description='Process some files.')
    parser.add_argument('--input', type=str, required=True, help='Path to the input JSON Shakespeare file')
    parser.add_argument('--output', type=str, required=True, help='Path to the output JSON Output file file')
    args = parser.parse_args()

    with open(args.output, 'r') as f:
        rawdata = json.load(f)
    with open(args.input, 'r') as f:
        trace = json.load(f)

    parsed_data = [
        (
            # data['metrics']['itl_ms'],
            # data['metrics']['e2e_ms'],
            # data['metrics']['ttft_ms'],
            # data['metrics']['input_tokens'],
            # data['metrics']['output_tokens'],
            pd.to_datetime(data['sent']),
            data['body'],
        )
        for data in rawdata['metrics']
    ]

    starttime = pd.to_datetime(rawdata['metadata']['start_time'])

    df = pd.DataFrame(
        parsed_data, 
        columns=[
            # 'itl_ms', 
            # 'e2e_ms', 
            # 'ttft_ms', 
            # 'input_tokens', 
            # 'output_tokens', 
            'sent', 
            'body'
        ]
    )

    parsed_trace = [
        (
            pd.to_timedelta(data['delay'], unit='ns'),
        )
        for data in trace
    ]

    df_trace = pd.DataFrame(
        parsed_trace, 
        columns=[
            'delay'
        ]
    )

    df_trace = df_trace.sort_values(by='delay')
    df = df.sort_values(by='sent')
    

    df['delay'] = (df['sent'] - starttime)
    df = df.reset_index(drop=True)
    df_trace = df_trace.reset_index(drop=True)

    df_merged = df['delay'] - df_trace['delay']
    df_merged_abs = df_merged.abs()

    min = df_merged_abs.dt.total_seconds().min()
    print(f"Min of df_merged: {min} seconds")
    p99 = np.percentile(df_merged_abs.dt.total_seconds(), 99)
    print(f"P99 of df_merged: {p99} seconds")

    plt.plot(df_merged.index, df_merged.dt.total_seconds(), label='Delay', marker='o')
    plt.grid(True)
    plt.xlabel('Time (s)')
    plt.ylabel('Delay (ns)')
    plt.title('Delay over Time')
    plt.show()

if __name__ == '__main__':
    main() 