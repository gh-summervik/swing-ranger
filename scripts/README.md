# Scripts

## collect_indices.sh

This script uses the `sr-cli` app to collect historical EOD prices for all the stocks in the SPY, QQQ, and DJIA. The stocks for which is searches are hard-coded, but current as of late May, 2026.

To run this script:

```bash
cd scripts/
chmod +x collect_indices.sh
./collect_indices.sh
```
