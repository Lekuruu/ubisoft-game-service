
import json
import os

def load_file(path: str) -> dict:
    with open(path, 'r') as file:
        return json.load(file)

def load_config(folder = './config') -> dict:
    config_files = [
        file.removesuffix('.json')
        for file in os.listdir(folder)
        if file.endswith('.json')
    ]

    return {
        file: load_file(f'{folder}/{file}.json')
        for file in config_files
    }
