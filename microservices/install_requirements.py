import pip
import sys


def install(packages):
    for package in packages:
        pip.main(['install', package])


def main():
    _all_ = [
        'grpcio==1.63.0',
        'grpcio-reflection==1.63.0',
        'matplotlib==3.7.5',
        'numpy==1.24.3',
        'pandas==2.0.3',
        'protobuf==5.26.1'
    ]
    windows = [
        'tensorflow==2.13.0'
    ]
    linux = [
        'tensorflow==2.13.0'
    ]
    darwin = [
        'tensorflow-metal==1.0.1',
        'tensorflow-macos==2.13.0'
    ]
    
    install(_all_)
    if sys.platform == 'windows':
        install(windows)
    elif sys.platform.startswith('linux'):
        install(linux)
    elif sys.platform == 'darwin':
        install(darwin)
    else:
        print(f'Error: Invalid platform ("{sys.platform}")', file=sys.stderr)

if __name__ == '__main__':
    main()
