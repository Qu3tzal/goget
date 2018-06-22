import uuid
import requests
import time
import sys
import multiprocessing
from threading import Thread

NUMBER_OF_REQUESTS_PER_TEST = 4

# Do one test, returns True if it succeded, False otherwise.
def do_one_test():
    id = str(uuid.uuid4())
    value = str(uuid.uuid4())
    url = 'http://localhost:8080/store/' + id
    data = '{"Value": "' + value + '"}'
    rtlist = [] # List of the response times

    start = time.perf_counter()
    response = requests.post(url, data=data)
    end = time.perf_counter()
    rtlist.append(end - start)

    if response.status_code != 200:
        print("Error on POST with UUID " + id + ". Status code : " + response.status_code)
        return False, 0

    start = time.perf_counter()
    response = requests.get(url, data="")
    end = time.perf_counter()
    rtlist.append(end - start)

    json_obj = response.json()
    if response.status_code != 200 or json_obj[id] != value:
        print("Error on GET with UUID " + id + ". Status code : " + response.status_code + ", value = '" + json_obj[id] + "' was expecting '" + value + "'")
        return False, 0

    start = time.perf_counter()
    response = requests.delete(url, data="")
    end = time.perf_counter()
    rtlist.append(end - start)

    if response.status_code != 200:
        print("Error on DELETE with UUID " + id + ". Status code : " + response.status_code)
        return False, 0

    start = time.perf_counter()
    response = requests.get(url, data="")
    end = time.perf_counter()
    rtlist.append(end - start)

    if response.status_code != 404:
        print("Error on GET with UUID after DELETE " + id + ". Status code : " + response.status_code)
        return False, 0

    mrt = 0
    for rt in rtlist:
        mrt += rt

    return True, mrt / NUMBER_OF_REQUESTS_PER_TEST

def do_multiple_tests(number_of_jobs):
    errors = 0
    mean_response_time_accu = 0
    start = time.perf_counter()

    for i in range(number_of_jobs):
        success, mean_response_time_for_test = do_one_test()

        if not success:
            errors += 1
        else: # We only take in account the mean response time for successful tests.
            mean_response_time_accu += mean_response_time_for_test

    end = time.perf_counter()
    mean_response_time = mean_response_time_accu / (number_of_jobs * NUMBER_OF_REQUESTS_PER_TEST)
    print(str(number_of_jobs) + " tests done in " + str(end - start) + "s. " + str(errors) + " failed. Mean response time = " + str(mean_response_time * 1000.0) + "ms.")

def main():
    threads = []

    number_of_jobs = 1000
    if len(sys.argv) > 1:
        number_of_jobs = int(sys.argv[1])

    number_of_threads = multiprocessing.cpu_count()
    jobs_per_thread = int(number_of_jobs / number_of_threads)

    print(str(number_of_threads) + " threads, " + str(jobs_per_thread) + " jobs per thread, 4 requests per job.")
    print("\t" + str(number_of_threads * jobs_per_thread * 4) + " requests")

    for i in range(number_of_threads):
        thread = Thread(target=do_multiple_tests, args=(jobs_per_thread,))
        thread.start()
        threads.append(thread)

    for t in threads:
        t.join()

    print("Tests finished.")

if __name__ == "__main__":
    main()
