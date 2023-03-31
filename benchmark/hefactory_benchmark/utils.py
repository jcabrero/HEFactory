import time, csv, os


# defining a decorator 
def test_case(func): 
    def wrapped(): 
    
        name = func.__name__
        
        func() 
        os.system("cp data/a.input.params ../lattigo_naive/bin/data/%s.input.params"%(name))
        os.system("cp data/a.input.params ../lattigo_performance/bin/data/%s.input.params"%(name))
    return wrapped 
    
class Timer(object):
    default_filename = "performance_overall.csv"
    def __init__(self, test_name, filename=None):
        self.times = dict()
        self.test_name = test_name
        self.filename = filename if not filename is None else Timer.default_filename

    def start(self, entry):
        self.times[entry] = time.time()
    
    def stop(self, entry):
        self.times[entry] =  time.time() - self.times[entry]
    
    def write_to_csv(self, row, header):
        with open(self.filename, 'a') as f:
            writer = csv.writer(f)
            # Write header if at start of file
            if f.tell() == 0:
                writer.writerow(header)
             # write a row to the csv file
            writer.writerow(row)

    def __del__(self):
        keys = list(self.times.keys())
        row = [self.test_name] + [self.times[k] for k in keys]
        header = ["Test#"] + keys
        self.write_to_csv(row, header)
