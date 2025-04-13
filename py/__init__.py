import time

def main_func():
    x = 1 + 1
    for i in range(20):
        print(i)
    print(x)
    print("You did it!!!")
    time.sleep(5)

    return x

if __name__ == '__main__':
    print(main_func())
