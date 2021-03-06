<coin-edit>
  <div class="modal is-active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">{ opts.coin ? "Edit " + opts.coin : "Add Coin" }</p>
        <button class="delete" aria-label="close" onclick={ close }></button>
      </header>
      <section class="modal-card-body">

        <div class="field" if={!opts.coin}>
            <label class="label">Coin</label>
            <div class="control">
              <select class="select" ref="coin">
                <option each={ coins } value={ short_name }>{ short_name }</option>
              </select>
            </div>
          </div>

        <div class="field">
          <label class="label">Pool URL</label>
          <div class="control">
            <input class="input" type="text" placeholder="stratum+tcp://example.com" ref="url" value={pool.url}>
          </div>
        </div>

        <div class="field">
          <label class="label">Pool User</label>
          <div class="control">
            <input class="input" type="text" placeholder="Usually your wallet address" ref="user" value={pool.user}>
          </div>
        </div>

        <div class="field">
          <label class="label">Pool Password</label>
          <div class="control">
            <input class="input" type="text" placeholder="Usually empty or x" ref="pass" value={pool.pass}>
          </div>
        </div>

        <article class="message is-danger" if={ validationError }>
          <div class="message-body">
            { validationError }
          </div>
        </article>

      </section>
      <footer class="modal-card-foot">
        <button class="button is-success" onclick={ save }>{ opts.coin ? "Save changes" : "Add coin" }</button>
        <button class="button" onclick={ close }>Cancel</button>
      </footer>
    </div>
  </div>

  <script>
    this.pool = opts.coin ? opts.miner.config.coins[opts.coin].pool : {}

    var configuredCoins = Object.keys(opts.miner.config.coins)
    this.coins = opts.miner.availableCoins.filter(function (available) {
      return !configuredCoins.find(function (configured) { return available.short_name == configured})
    })

    save() {
      var a = document.createElement('a')
      a.href = this.refs.url.value

      if (a.protocol != "stratum+tcp:" && a.protocol != "stratum+tls:") {
        this.validationError = "Pool URL must start with protocol (stratum+tcp:// or stratum+tls://)"
        return
      } else if (a.port == "") {
        this.validationError = "Pool URL must include port"
        return
      } else if (this.refs.user.value == "") {
        this.validationError = "Pool user must not be empty"
        return
      }

      opts.miner.config.coins[opts.coin || this.refs.coin.value] = {
        pool: {
          url: this.refs.url.value,
          user: this.refs.user.value,
          pass: this.refs.pass.value,
        }
      }
      opts.miner.trigger('update')
      opts.close()
    }

    close() {
      opts.close()
    }
  </script>
</coin-edit>
