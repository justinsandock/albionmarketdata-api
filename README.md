albiondata-api
==============

In short: this tool allows you to export prices over an api which you imported with albiondata-sql.

The API to access the data imported by [albiondata-sql](https://github.com/albiondata/albiondata-sql) from [albiondata-client](https://github.com/Regner/albiondata-client) over [albiondata-deduper](https://github.com/albiondata/albiondata-deduper/).


## Usage

Thanks to [viper](https://github.com/spf13/viper) and [cobra](https://github.com/spf13/cobra) you have 3 ways to configure albiondata-api.

### 1.) Traditional by configfile 

Just copy albiondata-api.yaml.tmpl to albiondata-api.yaml and edit it.

### 2.) By commandline arguments

See the output of the help page ```./albiondata-api -h```

### 3.) By environment variables

For example:

```
ADA_DBTYPE=sqlite3 ADA_DBURI=./sqlite.db ADA_LISTEN="[::]:3080" ./albiondata-api
```

## Authors

René Jochum <rene@jochums.at>


## LICENSE

MIT