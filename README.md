# mycache
Go exercise to implement a shared cache

download and go to mycache directory to compile

The program can run either by reading from a file - read and write actions to be applied to the cache or allowing the
program to create random read and write actions to be applied to the cache.

I recommend random actions.  If you choose random actions, then you can set an arbitrary number of threads and
read/write actions for each thread.

If you use files, you will need a file for every thread titled datafile\<thread number\>.dat with thread number starting
at 0.  For example datafile0.dat through datafile9 are included, more threads will require more files
  
You can also specify a file to initialize the cache, an initialization file is included, initData.dat

mycache options:

  -t \<number of threads as an integer\>
  
    specifies the number of threads
    
    
  -e \<number read/write actions as an integer\>
  
    specifies the number of read/write actions/events that each thread should execute
    
    
  -r
  
    uses the random event generator instead of files (the default is files)
    
    
  -m
  
    There are two types of cache controlling mechanisms.  No matter how many working threads, the default is to force
    all actions that hit the cache to serialize through a single thread.
    Using this option uses the multiplexor which allows multiiple threads to run in parallel.  Multiple readers run
    freely in parallel unless there is a writer.  Writers will block future readers and writers until its write is 
    complete, then freeing all other readers and writers.
    
    
  -i <initialization filename>
  
    Will read the specified file and use it to initialize the cache.  The format of the file is one row per key value
    pair to be initialized.  Each row has the following syntax
    
      <one character either i, s, or f> <integer key> <value>
      
      i specifies an integer value, s a string value, f a float value
      
      the integer key is an integer value that will be used as the key in the cache
      
      value is the value placed in the location key
