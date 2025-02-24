# 6Massive
6Massive is an efficient IPv6 large-scale target generation framework. Its idea is introduced in the paper "6Massive: An Efficient IPv6 Large-Scale Target Generation Framework". 6Massive can predict 1.644 billion IPv6 active addresses in 1.76 days.

## IPv6 Active Address List
The IPv6 active address probed by 6Massive is published on [IPv6 Active Address List](https://github.com/KeenoHao/IPv6_Active_Address_List.git).



## Execution steps

### Download seed address source Hitlist.

6Massive used the Hitlist of July 20, 2024. The latest Hitlist can be obtained in the following way.

   ```
   curl https://alcatraz.net.in.tum.de/ipv6-hitlist-service/open/responsive-addresses.txt.xz    --output responsive-addresses.txt.xz 
   xz -d responsive-addresses.txt.xz
   ```

### Compile 6Massive to generate an executable file.


   ```
   go mod init
   go mod tidy
   go build
   ```



### Construct seed sets.

Multiple seed address sets of the same scale can be obtained through expansion strategy.

The -S parameter is the seed address source, such as Hitlist. The -num parameter is the number of seed address sets. -Size is the size of the seed address set, that is, the number of seed addresses contained. -prefix is the prefix of the file name of the seed address set.

   ```
   ./6Massive -o expand -S hitlist_2024_07_20 -num 10 -size 100000 -prefix random100K
   ```


### IPv6 address pattern mining and IPv6 target address prediction.

Execute MDHC strategy on the input set of seed addresses, construct 4 IPv6 address space trees, and generate low-dimensional address patterns and high-dimensional address patterns.

The -s parameter represents the seed address set of the input, the -t parameter represents the target address generated in the low-dimensional pattern space, and the -h parameter represents the high-dimensional patterns.   

   ```
   ./6Massive -o MDHC -s random100K1 -t targetAddress -h highDimPattern
   ```

### Probe

Under a 30Mbps bandwidth, utilize the asynchronous scanning tool [ZMap](https://github.com/tumi8/zmap) to scan IPv6 target addresses and collect responsive IPv6 active addresses.

   ```
   sudo zmap --probe-module=icmp6_echoscan --ipv6-target-file=targetAddressFile  --output-file=activeAddressFile --ipv6-source-ip=(Machine IPv6 address) --bandwidth=30M --cooldown-time=4
   ```

### Feedback Strategy Based on Pattern Space Intersection.

1. Format active addresses.

The -a parameter represents the IPv6 active address probed in the low-dimensional pattern space, and the -t parameter represents the file output after formatting.
 
   ```
   ./6Massive -o convert -a activeAddress -t targetAddress
   ```

2. Execute feedback strategy.

The feedback strategy is used to filter out active high-dimensional patterns based on the active addresses in the low-dimensional pattern space, and generate IPv6 target addresses in these pattern spaces.
    
The -a parameter is the active address after formatting, -h is the high-dimensional address pattern file, -p is the pattern dimension of the high-dimensional address pattern, and -t is the target address to be output

   ```
   ./6Massive -o feedback -a activeAddress -h highDimPattern5 -p 5 -t targetAddress
   ```




## Reference

### Scanner

```
Durumeric, Z., Wustrow, E. & Halderman, J. A. ZMap: Fast Internet-wide scanning and its security applications. In 22nd
USENIX Security Symposium (USENIX Security 13), 605–620 (2013). https://www.usenix.org/conference/usenixsecurity13/
technical-sessions/paper/durumeric.

ZMap. ZMap Github code. https://github.com/tumi8/zmap.
```

### IPv6 Target Generation Algorithms

```
Murdock, A., Li, F., Bramsen, P., Durumeric, Z. & Paxson, V. Target generation for Internet-wide IPv6 scanning. In
Proceedings ofthe 2017 Internet Measurement Conference, 242–253, DOI: 10.1145/3131365.3131405 (2017).

Liu, Z., Xiong, Y., Liu, X., Xie, W. & Zhu, P. 6Tree: Efficient dynamic discovery of active addresses in the IPv6 address
space. Comput. Networks 155, 31–46, DOI: 10.1016/j.comnet.2019.03.010 (2019).

Hou, B., Cai, Z., Wu, K., Su, J. & Xiong, Y. 6Hit: A reinforcement learning-based approach to target generation for
Internet-wide IPv6 scanning. In IEEE INFOCOM 2021-IEEE Conference on Computer Communications, 1–10, DOI:
10.1109/INFOCOM42981.2021.9488794 (2021).

Yang, T., Cai, Z., Hou, B. & Zhou, T. 6Forest: An ensemble learning-based approach to target generation for Internet-
wide IPv6 scanning. In IEEE INFOCOM 2022-IEEE Conference on Computer Communications, 1679–1688, DOI:
10.1109/INFOCOM48880.2022.9796925 (2022).

Hou, B., Cai, Z., Wu, K., Yang, T. & Zhou, T. Search in the expanse: Towards active and global IPv6 hitlists. In IEEE
INFOCOM 2023-IEEE Conference on Computer Communications, 1–10, DOI: 10.1109/INFOCOM53939.2023.10229089
(2023).

Hou, B., Cai, Z., Wu, K., Yang, T. & Zhou, T. 6Scan: A high-efficiency dynamic Internet-wide IPv6 scanner with regional
encoding. IEEE/ACMTransactions on Netw. 31, 1870–1885, DOI: 10.1109/TNET.2023.3233953 (2023).

6Scan. 6Scan github code. https://github.com/hbn1987/6Scan.git.

Treestrace. Treestrace github code. https://github.com/6Seeks/Treestrace.git.
```

### Hitlists & Aliases

```
O. Gasser et al., “Clusters in the Expanse: Understanding and Unbiasing IPv6 Hitlists,” in IMC, 2018.

Hitlist. Hitlist data. https://alcatraz.net.in.tum.de/ipv6-hitlist-service/open/responsive-addresses.txt.xz

Aliases. Aliases data. https://alcatraz.net.in.tum.de/ipv6-hitlist-service/open/aliased-prefixes.txt.xz

Longest prefix matching for aliased prefixes. Longest prefix matching for aliased prefixes github code. https://ipv6hitlist.github.io/lpm/aliases-lpm.py
```

### Bloom Filter  & Deterministic Finite Automaton (DFA)

```
Bloom Filter. Bloom Filter github code. https://github.com/6Seeks/Treestrace.git.

DFA. DFA github code. https://github.com/NezhaFan/sieve.git.
```