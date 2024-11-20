# WorkerPool

This is a WorkerPool Implementation in golang with Auto Scaling.   

## Key Features  
 
1. One can set the max workers and min workers for auto scaling.   
2. If there are more jobs than workers, and less workers than max workers, new worker will be spawned to take up the job.   
3. If there are inactive workers and total workers is greater than min workers, inactive workers will be scaled down till min workers.   
4. Workers run concurrently and Isolated.   
5. For a incoming job, random available worker will be chosen.   

## Results   
 
| N_JOBS | JOB_TIME | N_GOROUTINES | N/2_WORKERS | 
| ------ | -------- | ------------ | ----------- |
| 1000000 |  500 ms  | 2.052139768 s   | 1.369614837 s |
| 500000| 500 ms |  1.358096717 s | 1.160396269 s   |
| 100000 | 500 ms | 681.6893570 ms | 1.039766089 s   |
| 1000000 | 700 ms | 2.344755774 s | 1.758263126 s |
| 500000 | 700 ms | 1.635180000 s | 1.582455652 s | 
| 100000 | 700 ms | 925.170906 ms | 1.44006322 s |

  
  
The above results are obtained with no auto scaling. The efficiency may be improved by enabling auto scaling for low amount of tasks and long running parrallel tasks.  
