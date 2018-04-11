# mycache
Go exercise to implement a shared cache

download and go to mycache directory to compile

the program can run from the included files to inject r and write actions to the cache, or it can be run with rand actions injected
I recommend random actions, then you can set the number of threads and number of read/write events/actions each thread will execute
if you use files you need a file for every thread titled datafile<thread number>.dat with thread numer starting at 0
For example datafile0.dat through datafile9 are included, more threads will requie more files
  
You can also specify a file to initialize the thread, and initialization file is included, initData.dat

mycache options:

  -t <number of threads as an integer>
  
    specifies the number of threads
    
    
  -e <number events as an integer>
  
    specifies the of read/write events events each thread should excute
    
    
  -r
  
    uses the random event generator instead of files (the default is files)
    
    
  -m
  
    There are two types of cache controlling mechanisms.  The default is to force all actions that hit the cache through a single
    thread.
    Using this option uses the multiplexor which allows multiiple threads to run in parallel.  Multiple readers run freely in parallel
    unless there is a writer.  Writers will block future readers and writers until its write is complete, then freeing other readers
    and writers.
    
    
  -i <initialization filename>
  
    Will read the specified file and use it to initialize the cache.  The format of the file is one row per key value pair to be 
    initialized.  Each row has the following syntax
    
      <one character either i, s, or f> <integer key> <value?
      
      i specifies an integer value, s a string value, f a float value
      
      the integer key is an integer value that will be used as the key in the cache
      
      value is the value placed in the location key
