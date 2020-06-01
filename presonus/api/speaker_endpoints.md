### Existing Endpoints
* "Speaker/75hz": "0..1",  //100Hz high pass filter on/off
* "Speaker/contour": "0/0.5/1", //0: Normal, 0.5: LBR Source, 1: Floor Monitor
* "Speaker/clip": "0/1", //speaker digital clip LED flash
* "Speaker/limit": "0/1", //Speaker limit LED flash?
* "Speaker/line/ch1/delay": "0..1", //Actual range is 0..500 ms
* "Speaker/line/ch1/delay_enable": "0/1",
* "Speaker/line/ch1/eq/eqallon": "0/1", //Param EQ Enable
* "Speaker/line/ch1/eq/eqbandon1": "0/1", //Param EQ1 Enable
* "Speaker/line/ch1/eq/eqbandon8": "0/1", //Param EQ8 Enable
* "Speaker/line/ch1/eq/eqfreq1": "0..1", //Param EQ1 Freq.  Actual range is 20..20,000 Hz
* "Speaker/line/ch1/eq/eqfreq8": "0..1", //Param EQ8 Freq.  Actual range is 20..20,000 Hz
* "Speaker/line/ch1/eq/eqgain1": "0..1", //Param EQ1 Gain.  Actual range is -15..15 dB
* "Speaker/line/ch1/eq/eqgain8": "0..1", //Param EQ8 Gain.  Actual range is -15..15 dB
* "Speaker/line/ch1/eq/eqq1": "0..1", //Actual range is 0.10..4
* "Speaker/line/ch1/eq/eqq8": "0..1", //Actual range is 0.10..4
* "Speaker/line/ch1/geq/on": "0/1", //Graphic EQ Enable
* "Speaker/line/ch1/geq/gain1": "0..1",  //Graphic EQ01 Gain.  Actual range is -16..16 dB
* "Speaker/line/ch1/geq/gain31": "0..1", //Graphic EQ31 Gain.  Actual range is -16..16 dB
* "Speaker/line/ch1/limit/limiteron": "0/1",
* "Speaker/line/ch1/limit/threshold": "0..1", //actual range is -28..0 dB
* "Speaker/line/ch1/mute": "0/1",
* "Speaker/line/ch1/notch/notchallon": "0/1", //Notch EQ Enable
* "Speaker/line/ch1/notch/notchbandon1": "0/1", //Notch EQ1 Enable
* "Speaker/line/ch1/notch/notchbandon8": "0/1", //Notch EQ8 Enable
* "Speaker/line/ch1/notch/notchfreq1": "0..1", //Notch EQ1 Freq.  Actual range is 20..20,000 Hz
* "Speaker/line/ch1/notch/notchfreq8": "0..1", //Notch EQ8 Freq.  Actual range is 20..20,000 Hz
* "Speaker/line/ch1/notch/notchgain1": "0..1", //Notch EQ1 Gain.  Actual range is -48..0 dB
* "Speaker/line/ch1/notch/notchgain8": "0..1", //Notch EQ8 Gain.  Actual range is -48..0 dB
* "Speaker/line/ch1/volume": "0..1", //Actual range is -84..10 dB
* "Speaker/line/ch1/volume_enable": "0/1",
* "Speaker/presetloading": "0", //seems to show while the speaker is muted to change modes
* "Speaker/remote_lim_on": "0/1", //follows Speaker/line/ch1/limit/limiteron
* "Speaker/remote_notch_on": "0/1", //follows Speaker/line/ch1/notch/notchallon
* "Speaker/remote_geq_on": "0/1", //follows Speaker/line/ch1/geq/on
* "Speaker/remote_peq_on": "0/1", //follows Speaker/line/ch1/eq/eqallon
* "Speaker/remotelayeron": "0/1", //User mode Enable
* "Speaker/signal": "0/1", //Speaker signal LED flash
* "Speaker/wink": "0/1", //Front LED Color.  0: Blue, 1:White
* "Speaker/excursion": "0/1", //Speaker excursion warning

### Missing but expected:
* Something about amp temp
* Something about digital peak 0/1
* Something about limiter hitting 0/1
* something about temp too high 0/1  (maybe covered by amp temp)


### Parametric & Notch EQ defaults
* Speaker/line/ch1/eq && Speaker/line/ch1/notch endpoints have settings for bands 1..8

Default Frequencies Bands:
1. 60Hz
2. 120Hz
3. 240Hz
4. 480Hz
5. 960Hz
6. 1.92kHz
7. 3.84kHz
8. 7.68kHz

Note: Notch EQs have no Q parameter

### Graphic EQ
* Speaker/line/ch1/geq/gain endpoints have levels for bands 1..31

Frequencies:
* Band 01: 20
* Band 02: 25
* Band 03: 32
* Band 04: 40
* Band 05: 50
* Band 06: 63
* Band 07: 80
* Band 08: 100
* Band 09: 125
* Band 10: 160
* Band 11: 200
* Band 12: 250
* Band 13: 320
* Band 14: 400
* Band 15: 500
* Band 16: 640
* Band 17: 800
* Band 18: 1k
* Band 19: 1.3k
* Band 20: 1.6k
* Band 21: 2k
* Band 22: 2.6k
* Band 23: 3.2k
* Band 24: 4k
* Band 25: 5.2k
* Band 26: 6.4k
* Band 27: 8k
* Band 28: 10k
* Band 29: 13k
* Band 30: 16k
* Band 31: 20k


### Debugging

DEBUGGING by trying to connect WORX tool to speaker, unsuccesfully.  But found 'digitalin' endpoint for dante control

* Speaker/monosum.......UC..$.PVd.f.
* Speaker/remote_notch_on.......UC....PVd.f.
* Speaker/line/ch1/notch/notchallon.......UC..".PVd.f." +
* Speaker/remote_geq_on.......UC..$.PVd.f.
* Speaker/line/ch1/geq/on.......UC....PVd.f.
* Speaker/digitalin......?UC..$.PVd.f.
* Speaker/line/ch1/volume....A5?UC..+.PVd.f.
* Speaker/line/ch1/volume_enable.......UC..#.PVd.f.
* Speaker/line/ch1/delay.......UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype1.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq1...&nH?UC..(.PVd.f.
* Speaker/line/ch1/eq/eqfreq1...}.">UC..*.PVd.f." +
* Speaker/line/ch1/eq/eqbandop1.......UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype2.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq2......?UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype3.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq3......?UC..(.PVd.f.
* Speaker/line/ch1/eq/eqfreq3...%..>UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype4.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq4......?UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype5.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq5......?UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype6.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq6......?UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype7.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq7......?UC..(.PVd.f.
* Speaker/line/ch1/eq/eqtype8.......UC..%.PVd.f.
* Speaker/line/ch1/eq/eqq8...&nH?UC..*.PVd.f.
* Speaker/line/ch1/eq/eqbandop8.......UC....PVd.f.
* Speaker/line/ch1/notch/notchfreq1...}.">UC....PVd.f.
* Speaker/line/ch1/notch/notchfreq3...%..>UC....KAf.d.UC....PVd.f.
* Speaker/wink.
* Speaker/line/ch1/eq/eqfreq1...z.">UC..(.PVd.f.
* Speaker/line/ch1/eq/eqfreq3...$..>UC....PVd.f.
* Speaker/line/ch1/notch/notchfreq1...z.">UC....PVd.f.
* Speaker/line/ch1/notch/notchfreq3
