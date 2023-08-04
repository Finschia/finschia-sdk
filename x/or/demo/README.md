# Demo

## Settlement

Run a chain and challenge game.

```bash
$ ./run_chain.sh
$ ./challenge.sh
```

Example logs are below. It is ok if `demo successful` is output.

```bash
$ make install # install simd
$ cd x/or/demo/tools
$ make all     # build tools
$ cd ..
$ ./run_chain.sh
$ ./challenge.sh
# Start challenge
http://localhost:26657/tx?hash=0x732DE03CD38EC3934DF523BEDE477DE155B039330DA460C8A66239DFE1650F72
# Nsect challenge
## Round 1
http://localhost:26657/tx?hash=0x35AE5F6680D4B5072449F2FC24578BB7D7F6CBB07DB02CF76D9A4D8398B54E2D
http://localhost:26657/tx?hash=0x35BB8652D4AB58606BF5A6B44C0519FACC95258BEC9BA714640593BB650F6494
## Round 2
http://localhost:26657/tx?hash=0xCCC6F23815D33EE518E48909CB666718B13F75F163BFF71FB892F799EC2FD31A
http://localhost:26657/tx?hash=0x1EA1FF8475CD7E5825620DFDA775DBC3BDADBADBA96EC0ABBE8DFAB6F90BAF2C
## Round 3
http://localhost:26657/tx?hash=0x514EFB778AC929E681AF5534E39A825A8F1CCDFC6987B77FF206FAAE1232E158
http://localhost:26657/tx?hash=0x966EF886F0A70B1A3036F7D6BA3506F50B38F53FA50C45160D90C78C8E2CD585
## Round 4
http://localhost:26657/tx?hash=0x5914952B84BCF32CF1EE20AB8B65250A0132AA614114958A15C1700CBD613FFE
http://localhost:26657/tx?hash=0x578F2DCBB87224EEF10FD1952A9254A7BA6BEEFF31490ECBA252D6114F64B46D
## Round 5
http://localhost:26657/tx?hash=0x2E879A29D0C3838B2BE5791473A70F1817D6D07289070E6EC473079133206550
http://localhost:26657/tx?hash=0xDCFC27045A959456D017DDB831FE81E2BDA24EF816F3E1684BB2DA643FCC69A2
# Finish challenge
http://localhost:26657/tx?hash=0xEBADA91B858A79F2F7A8AE587F6BBAAEF20097C527FAA7FAB80FCC7B24ED1693
Challenger wins:) demo successful:)
```

## Tools

- cannon - Mips vm, execute step, generate state and witness.
- challenge - Challenge structure emulator, genereate challenge id, steps, and flag if challenge is searching or not.
- miniapp - Verification program, verify if sequecner's state is correct or not.
- preimage - Preimage oracle server, provide miniapp preimage.
    - data/height-100.json - Correct data for height 100.
    - data/height-100-def.json - Malicious data. This demo assume that this data is proposed by malicious sequencer.
    - data/height-100-cha.json - Data for verification. Challenger prove sequencer's fraud using this data.
