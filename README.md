<div id="top"></div>
<!-- PROJECT LOGO -->
<br />
<div align="center">

  <img src="./.github/assets/topos_logo.png#gh-light-mode-only" alt="Logo" width="200">
  <img src="./.github/assets/topos_logo_dark.png#gh-dark-mode-only" alt="Logo" width="200">

<br />

<p align="center">
<em>go-topos-sequencer-client</em> is an early proof of concept client for Golang GRPC bidirectional communication.
</p>

<br />

</div>

## Background
The Topos Sequencer produces and signs certificates using [ICE-FROST](https://eprint.iacr.org/2021/1658.pdf) signatures. It requires communication between sequencers of the same subnet to proceed to the DKG to generate the group key, and sign certificates with joint effort using the threshold signature protocol.

One way to implement this protocol is to reuse the existing p2p network that the subnet nodes maintain. `go-topos-sequencer-client` will eventually be used for bidirectional communication between subnet nodes and their Topos sequencers to generate the group key and sign certificates by leveraging the subnet p2p network.


## License

This project is released under the terms specified in the [LICENSE](LICENSE) file.
