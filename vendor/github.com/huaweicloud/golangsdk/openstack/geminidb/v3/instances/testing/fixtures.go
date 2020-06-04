package testing

const CreateRequest = `
{ 
  "name": "test-cassandra-01", 
  "datastore": { 
    "type": "GeminiDB-Cassandra", 
    "version": "3.11", 
    "storage_engine": "rocksDB" 
  }, 
  "region": "aaa", 
  "availability_zone": "bbb", 
  "vpc_id": "674e9b42-cd8d-4d25-a2e6-5abcc565b961", 
  "subnet_id": "f1df08c5-71d1-406a-aff0-de435a51007b", 
  "security_group_id": "7aa51dbf-5b63-40db-9724-dad3c4828b58", 
  "password": "Test@123", 
  "mode": "Cluster", 
  "flavor": [ 
    { 
      "num": 3, 
      "size": 500,
      "storage": "ULTRAHIGH",
      "spec_code": "nosql.cassandra.4xlarge.4" 
    } 
  ], 
  "backup_strategy": { 
    "start_time": "08:15-09:15", 
    "keep_days": 8
  },
  "enterprise_project_id": "0" 
}     `

const CreateResponse = `
{ 
  "id": "39b6a1a278844ac48119d86512e0000bin06", 
  "name": "test-cassandra-01", 
  "datastore": { 
    "type": "GeminiDB-Cassandra", 
    "version": "3.11", 
    "storage_engine": "rocksDB" 
  }, 
  "created": "2019-10-28 14:10:54",
  "status": "creating",
  "region": "aaa", 
  "availability_zone": "bbb", 
  "vpc_id": "674e9b42-cd8d-4d25-a2e6-5abcc565b961", 
  "subnet_id": "f1df08c5-71d1-406a-aff0-de435a51007b", 
  "security_group_id": "7aa51dbf-5b63-40db-9724-dad3c4828b58", 
  "mode": "Cluster", 
  "flavor": [ 
    { 
      "num": 3, 
      "size": 500,
      "storage": "ULTRAHIGH",
      "spec_code": "nosql.cassandra.4xlarge.4" 
    } 
  ], 
  "backup_strategy": { 
    "start_time": "08:15-09:15", 
    "keep_days": 8 
  } ,
  "job_id": "c010abd0-48cf-4fa8-8cbc-090f093eaa2f",
  "enterprise_project_id": "0" 
}
    `
