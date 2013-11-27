package ssm

const INIT_SYSTEMD byte = 0x0

const CMD_SYSTEMCTL string = "/sbin/systemctl"

const SERVICE_SERVICE_RUNNING       string = " is running as pid "
const SERVICE_SERVICE_NOT_RUNNING   string = " is not running"
const SERVICE_SERVICE_NOT_ENABLED   string = "_enable to YES in"
