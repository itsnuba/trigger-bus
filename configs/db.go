package configs

type configDb struct {
	MongoUri                 string `mapstructure:"MONGO_URI"`
	MongoDB                  string `mapstructure:"MONGO_DB"`
	MongoColEvent            string `mapstructure:"MONGO_COLLECTION_EVENT"`
	MongoColTriggerListener  string `mapstructure:"MONGO_COLLECTION_TRIGGER_LISTENER"`
	MongoColTriggerScheduler string `mapstructure:"MONGO_COLLECTION_TRIGGER_SCHEDULER"`
	MongoColTriggerLog       string `mapstructure:"MONGO_COLLECTION_TRIGGER_LOG"`
}
