# WorkerPool

This is yet another WorkerPool Implementation with Auto Scaling.   

## Key Features  
 
1. One can set the max workers and min workers for auto scaling.  
2. If there are more jobs than workers, and less workers than max workers, new worker will be spawned to take up the job.  
3. If there are inactive workers and total workers is greater than min workers, inactive workers will be scaled down till min workers.  
4. Workers run concurrently and Isolated.
5. For a incoming job, random available worker will be chosen.  
