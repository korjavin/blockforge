<gui-main>
  <nav class="navbar is-primary" role="navigation" aria-label="main navigation">
    <div class="container">
      <div class="navbar-brand">
        <div class="navbar-brand">
          <a class="navbar-item has-text-weight-bold is-size-5">
            BlockForge
          </a>
        </div>
      </div>
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="field is-grouped">
            <p class="control">
              <a class={"button is-large " + (miner.running ? "is-warning" : "is-success") } onclick={toggleRunning}>
                <span class="icon">
                  <i class={"fa " +(miner.running ? "fa-pause" : "fa-play")}></i>
                </span>
              <span>
                  { miner.running ? "Stop mining" : "Start mining"}
              </span>
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </nav>

  <div class="container">
    <section class="section">
      <hardware-list miner={miner}></cpu-list>
    </section>
    <section class="section" style="padding-top: 0">
      <coin-list miner={miner}></coin-list>
    </section>
  </div>
  <script>

  window.miner = this.miner = new Miner()

  this.miner.on('updated', function() {
    this.update()
  }.bind(this))

  toggleRunning() {
    if (this.miner.running) {
      this.miner.trigger('stop')
    } else {
      this.miner.trigger('start')
    }
    this.update()
  }

  function Miner() {
    riot.observable(this)

    var backend = opts.backend
    var updateInterval

    this.config = backend.data.config
    this.availableCoins = backend.data.coins
    this.processors = backend.data.processors
    this.running = false

    this.on('start', function() {
      backend.start()
      this.running = true
      updateInterval = setInterval(function() {
        backend.stats()
      }, 1000)
    })

    this.on('stop', function() {
      clearInterval(updateInterval)
      backend.stop()
      this.running = false
    })

    this.on('update', function() {
      backend.updateConfig(JSON.stringify(this.config))
      this.trigger('updated')
    })
  }
  </script>
</gui-main>
