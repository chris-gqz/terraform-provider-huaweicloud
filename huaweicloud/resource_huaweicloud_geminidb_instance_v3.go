package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/instances"
)

var projectID string

func resourceGeminiDBInstanceV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeminiDBInstanceV3Create,
		Read:   resourceGeminiDBInstanceV3Read,
		Update: resourceGeminiDBV3Update,
		Delete: resourceGeminiDBInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"db_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "GeminiDB-Cassandra",
							ValidateFunc: validation.StringInSlice([]string{
								"GeminiDB-Cassandra",
							}, true),
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "rocksDB",
							ValidateFunc: validation.StringInSlice([]string{
								"rocksDB",
							}, true),
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "3.11",
							ValidateFunc: validation.StringInSlice([]string{
								"3.11",
							}, true),
						},
					},
				},
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceGeminiDBDataStore(d *schema.ResourceData) instances.DataStore {
	var dataStore instances.DataStore
	datastoreRaw := d.Get("datastore").([]interface{})
	log.Printf("[DEBUG] datastoreRaw: %+v", datastoreRaw)
	if len(datastoreRaw) == 1 {
		dataStore.Type = datastoreRaw[0].(map[string]interface{})["engine"].(string)
		dataStore.Version = datastoreRaw[0].(map[string]interface{})["version"].(string)
		dataStore.StorageEngine = datastoreRaw[0].(map[string]interface{})["storage_engine"].(string)
	} else {
		dataStore.Type = "GeminiDB-Cassandra"
		dataStore.Version = "3.11"
		dataStore.StorageEngine = "rocksDB"
	}
	log.Printf("[DEBUG] datastore: %+v", dataStore)
	return dataStore
}

func resourceGeminiDBBackupStrategy(d *schema.ResourceData) instances.BackupStrategy {
	var backupStrategy instances.BackupStrategy
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	log.Printf("[DEBUG] backupStrategyRaw: %+v", backupStrategyRaw)
	if len(backupStrategyRaw) == 1 {
		backupStrategy.StartTime = backupStrategyRaw[0].(map[string]interface{})["start_time"].(string)
		backupStrategy.KeepDays = strconv.Itoa(backupStrategyRaw[0].(map[string]interface{})["keep_days"].(int))
	} else {
		backupStrategy.StartTime = "00:00-01:00"
		backupStrategy.KeepDays = "7"
	}
	log.Printf("[DEBUG] backupStrategy: %+v", backupStrategy)
	return backupStrategy
}

func resourceGeminiDBFlavor(d *schema.ResourceData) []instances.Flavor {
	var flavorList []instances.Flavor
	flavor := instances.Flavor{
		Num:      "3",
		Size:     strconv.Itoa(d.Get("volume_size").(int)),
		Storage:  "ULTRAHIGH",
		SpecCode: d.Get("flavor").(string),
	}
	flavorList = append(flavorList, flavor)
	return flavorList
}

func GeminiDBInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := instances.ListGeminiDBInstanceOpts{
			Id: instanceID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return nil, "", err
		}
		instancesList, err := instances.ExtractGeminiDBInstances(allPages)
		if err != nil {
			return nil, "", err
		}

		if instancesList.TotalCount == 0 {
			var instance instances.GeminiDBInstanceResponse
			return instance, "deleted", nil
		}
		insts := instancesList.Instances

		return insts[0], insts[0].Status, nil
	}
}

func resourceGeminiDBInstanceV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s ", err)
	}

	createOpts := instances.CreateGeminiDBOpts{
		Name:             d.Get("name").(string),
		DataStore:        resourceGeminiDBDataStore(d),
		Region:           GetRegion(d, config),
		AvailabilityZone: d.Get("availability_zone").(string),
		VpcId:            d.Get("vpc_id").(string),
		SubnetId:         d.Get("subnet_id").(string),
		SecurityGroupId:  d.Get("security_group_id").(string),
		Password:         d.Get("password").(string),
		Mode:             "Cluster",
		BackupStrategy:   resourceGeminiDBBackupStrategy(d),
		Flavor:           resourceGeminiDBFlavor(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting instance from result: %s ", err)
	}
	log.Printf("[DEBUG] Create : instance %s: %#v", instance.Id, instance)

	d.SetId(instance.Id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"normal"},
		Refresh:    GeminiDBInstanceStateRefreshFunc(client, instance.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      120 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		errDelete := resourceGeminiDBInstanceV3Delete(d, meta)
		if errDelete != nil {
			return fmt.Errorf("Error delete geminidb fail err: %s ", err)
		}
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instance.Id, err)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandGeminiDBTags(tagRaw)
		projectID = client.ProjectID
		client.ProjectID = ""
		client.ResourceBase = strings.TrimRight(client.ResourceBase, "/")
		if tagErr := tags.Create(client, "instances", instance.Id, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of GeminiDB %q: %s", instance.Id, tagErr)
		}
	}

	return resourceGeminiDBInstanceV3Read(d, meta)
}

func resourceGeminiDBInstanceV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}

	instanceID := d.Id()
	opts := instances.ListGeminiDBInstanceOpts{
		Id: instanceID,
	}
	client.ProjectID = projectID
	client.ResourceBase = strings.TrimRight(client.ResourceBase, "/") + "/" + projectID
	if !strings.HasSuffix(client.ResourceBase, "/") {
		client.ResourceBase = client.ResourceBase + "/"
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return fmt.Errorf("Error fetching GeminiDB instance: %s", err)
	}
	instances, err := instances.ExtractGeminiDBInstances(allPages)
	if err != nil {
		return fmt.Errorf("Error extracting GeminiDB instance: %s", err)
	}
	if instances.TotalCount == 0 {
		return fmt.Errorf("Error fetching GeminiDB instance: deleted")
	}
	insts := instances.Instances
	instance := insts[0]

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("availability_zone", instance.AvailabilityZone)
	d.Set("name", instance.Name)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("mode", instance.Mode)
	d.Set("db_username", instance.DbUserName)
	d.Set("status", instance.Status)
	d.Set("port", instance.Port)

	dbList := make([]map[string]interface{}, 0, 1)
	db := map[string]interface{}{
		"engine":         "GeminiDB-Cassandra",
		"version":        "3.11",
		"storage_engine": "rocksDB",
	}
	dbList = append(dbList, db)
	d.Set("datastore", dbList)

	nodesList := make([]map[string]interface{}, 0, 1)
	if len(instance.Groups) > 0 {
		for _, group := range instance.Groups {
			for _, Node := range group.Nodes {
				node := map[string]interface{}{
					"id":         Node.Id,
					"name":       Node.Name,
					"status":     Node.Status,
					"private_ip": Node.PrivateIp,
				}
				nodesList = append(nodesList, node)
			}
		}
	}
	d.Set("nodes", nodesList)

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	d.Set("backup_strategy", backupStrategyList)

	//save geminidb tags
	client.ProjectID = ""
	client.ResourceBase = strings.TrimRight(client.ResourceBase, "/")
	resourceTags, err := tags.Get(client, "instances", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud geminidb tags: %s", err)
	}

	tagmap := make(map[string]string)
	for _, val := range resourceTags.Tags {
		tagmap[val.Key] = val.Value
	}
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags for HuaweiCloud geminidb (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceGeminiDBInstanceV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s ", err)
	}

	instanceId := d.Id()
	result := instances.Delete(client, instanceId)
	if result.Err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal", "abnormal", "creating", "createfail", "enlargefail", "data_disk_full"},
		Target:     []string{"deleted"},
		Refresh:    GeminiDBInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	return nil
}

func expandGeminiDBTags(tagmap map[string]interface{}) []tags.ResourceTag {
	var taglist []tags.ResourceTag

	for k, v := range tagmap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return taglist
}

func resourceGeminiDBV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc: %s", err)
	}
	//update tags
	if d.HasChange("tags") {
		//remove old tags and set new tags
		old, new := d.GetChange("tags")
		oldRaw := old.(map[string]interface{})
		if len(oldRaw) > 0 {
			taglist := expandGeminiDBTags(oldRaw)
			if tagErr := tags.Delete(client, "geminidb", d.Id(), taglist).ExtractErr(); tagErr != nil {
				return fmt.Errorf("Error deleting tags of GeminiDB %q: %s", d.Id(), tagErr)
			}
		}

		newRaw := new.(map[string]interface{})
		if len(newRaw) > 0 {
			taglist := expandGeminiDBTags(newRaw)
			client.ProjectID = ""
			client.ResourceBase = strings.TrimRight(client.ResourceBase, "/")
			if tagErr := tags.Create(client, "instances", d.Id(), taglist).ExtractErr(); tagErr != nil {
				return fmt.Errorf("Error setting tags of GeminiDB %q: %s", d.Id(), tagErr)
			}
		}
	}

	return resourceGeminiDBInstanceV3Read(d, meta)
}
