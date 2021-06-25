## What is this?

This is basic PoC showing how minipool validator keys could be split among oDAO members, through the use of threshold cryptography.

## How do I run it?

```sh
docker build -t rocketpool-split-keys .
docker run --rm -ti rocketpool-split-keys
```

## Sample run output

```
Created minipool. Address: 0xdeadbeef, Validator pubkey: 1 1250ddbb2aa5d57ae756540362ad4efb4257ebb50add28d80a7bbc3f6d9ec3[...]
Sending key shares to the ODAO...
Sending minipool '0xdeadbeef' share to ODAO member '1': 317fdaffc3a77608405aeae42016477138f7a9501e1e806e1370365a17f4904d
Sending minipool '0xdeadbeef' share to ODAO member '2': ad4c0033ad87b79670885541b198cbf7c30ee8c5643d33c5adf54f7287644a8
Sending minipool '0xdeadbeef' share to ODAO member '3': 4f7440756fc6db7eef18a14f0b6938b421f491e9dcfc1cdbe7e2708fac946e56
Sending minipool '0xdeadbeef' share to ODAO member '4': 17830db00f379b8872178ec4ddc19b4482c74b62b24aa54eba798925a44f0d55
Sending minipool '0xdeadbeef' share to ODAO member '5': 4adc76596c65b6265678fdc5a566647b462462fcd62c2492d2a49eb70fa621a7
Sending minipool '0xdeadbeef' share to ODAO member '6': 1a52bcb341630c835c93e414f13e44dc49090b248a3e2aa3063b145ee99ab4a
Sending minipool '0xdeadbeef' share to ODAO member '7': 23b87cabb98405fe767c0047ee0dcac6a5871c8909ae9792d3b6c0d04129aa40
Sending minipool '0xdeadbeef' share to ODAO member '8': 3d28c1a7d311b880e5576bd178b23fe0954a627e194de74dbc9dcd5707561e88
Sending minipool '0xdeadbeef' share to ODAO member '9': 4df5fabf80bf484f825b80ddef01439b93da62917781d1daeb18d6da411f0822
Sending minipool '0xdeadbeef' share to ODAO member '10': 562027f2c28cb56a4d883f6d50fad5f7a1371cc3244a573a5f27dd59ee84670e
Will try to recover validator key using 3 shares...
Using ODAO member  8 key share: 3d28c1a7d311b880e5576bd178b23fe0954a627e194de74dbc9dcd5707561e88
Using ODAO member  1 key share: 317fdaffc3a77608405aeae42016477138f7a9501e1e806e1370365a17f4904d
Using ODAO member  5 key share: 4adc76596c65b6265678fdc5a566647b462462fcd62c2492d2a49eb70fa621a7
Successfully verified key shares.
Recovered key:  1 1250ddbb2aa5d57ae756540362ad4efb4257ebb50add28d80a7bbc3f6d9ec3[...]
Contract key:   1 1250ddbb2aa5d57ae756540362ad4efb4257ebb50add28d80a7bbc3f6d9ec3[...]
Sucessfully started staking on minipool '0xdeadbeef'...
```

