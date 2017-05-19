#!/usr/bin/env bash

/usr/bin/time -f "\nCPU: %Us\tReal: %es\tRAM: %MKB" go run pg.go lib.go
/usr/bin/time -f "\nCPU: %Us\tReal: %es\tRAM: %MKB" go run pg-jsonb.go lib.go
/usr/bin/time -f "\nCPU: %Us\tReal: %es\tRAM: %MKB" go run cockroach.go lib.go
/usr/bin/time -f "\nCPU: %Us\tReal: %es\tRAM: %MKB" go run scylla.go lib.go
