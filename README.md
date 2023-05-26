<p align="center">
  <img src="pocsagpilemaster.png" />
</p>

# POCSAG PILEMASTER
A filter for rtl_sdr and multimon-ng that sends POCSAGS message into an MQTT broker.
Only using it on a Raspberry Pi Zero 2 W Rev 1.0. 

Might add TLS support one day, hopefully the zero can handle it.


# Example usage


Requires rtl_fm and multimon-ng and a compatible rtl-sdr.

This example also saves all the messages to file (unfiltered - that you can replay), while pocsagpilemaster does some filtering.


```
export POCSAGPILEMASTER_BROKER=mahbrokah.evilmega.corp
export POCSAGPILEMASTER_PORT=1883
export POCSAGPILEMASTER_CLIENTID=pocsag0001
export POCSAGPILEMASTER_USERNAME=pocsag
export POCSAGPILEMASTER_PASSWORD=
export POCSAGPILEMASTER_TOPIC="mypocsags/pocsag01"
export POCSAGPILEMASTER_DEBUG=YES

rtl_fm -M fm -f 169.800M -s 22050 -g 100 -l 310 | multimon-ng --timestamp -e -u -C SE -t raw -a POCSAG512 -a POCSAG1200 -a POCSAG2400 -f alpha -a scope /dev/stdin | tee -a pocsag.out | pocsagpilemaster 
```

# Example build
Crosscompiling:

```
env GOOS=linux GOARCH=arm64 go build
```
