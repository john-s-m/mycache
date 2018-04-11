# mycache
Go exercise to implement a shared cache

download and go to mycache directory to compile

The program can run eiter by reading read and write actions from files or creating random read and write actions to the cache.


I recommend random actions.  If you choose random actions, then you can set an arbitrary number of threads and read/write
events/actions for each thread.

If you use files, you will need a file for every thread titled datafile<thread number>.dat with thread number starting at 0
For example datafile0.dat through datafile9 are included, more threads will require more files
  
You can also specify a file to initialize the cache, an initialization file is included, initData.dat

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
    
      <one character either i, s, or f> <integer key> <value>
      
      i specifies an integer value, s a string value, f a float value
      
      the integer key is an integer value that will be used as the key in the cache
      
      value is the value placed in the location key
