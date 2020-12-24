# API 

The end user API is designed to run a long running job from the server. We can also control a long running upload from the client as well.

Once we start the request, we can track the progress of the request through api endpoints, and pause, terminate the task as well. This makes sure we have full control over the long running job.

## API Endpoints
- General Server API - We can limit the number of concurrent jobs that can be running in the server. Get/Set the current running task with a worker id.
	* `new_job/` `GET` - Check if a new job can be concurrently handled. 
	  * Return a **200** response if a slot is available, along with a `worker-id` for making a request. `worker-id` expires in 5 minutes, and should be used within that period. 
	  * Return a **500** if server error.
	* `modify/{worker-id}/{new_job_status}` `POST` - Make a POST request for modifying the state of the task. 
	  * Return **200** with new status of the job. Refer [this]() for possible states from the given initial state.
	  * Return **400** if the job status is not correct, or the worker id is invalid.
	* `status/{worker-id}` `GET` - Get the current status of the running task.
	  * Return **200** with status of the job.
	  * Return **400** if the worker-id is invalid.
- Uploading Data - 
	- `upload/new/{worker-id}` `POST`- Make a POST request for uploading data.
	  - Return **200** if the task is terminated/completed. Return the final state of the request if completed, or was terminated halfway.
	  - Return a **503** if the server has no available slot. 
	  - Return **400** if the worker id is invalid.
	  - Return **500** if the request cannot be completed.
- Downloading Data - 
  - `download/new/{worker-id}` `POST` - Make a POST request for downloading data from server.
    - Return **200** if the task is terminated/completed. Return the final state of the request if completed, or was terminated halfway.
    - Return a **503** if the server has no available slot. 
    - Return **400** if the worker id is invalid.
    - Return **500** if the request cannot be completed.
- Running a long running job - 
  - `job/new/{worker-id}` `POST` - Initialise a long running job which can be server processing.
    - Return **200** if the task is terminated/completed. Return the final state of the request if completed, or was terminated halfway.
    - Return a **503** if the server has no available slot. 
    - Return **400** if the worker id is invalid.
    - Return **500** if the request cannot be completed.



