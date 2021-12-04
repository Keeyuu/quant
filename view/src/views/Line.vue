<template>

<div
    class="k-line-chart-container">
    <div id=choice>
      <input type="text" v-model="code" placeholder="edit me">
      <input type="text" v-model="level" placeholder="edit me">
      <button
        v-on:click="select_code()">
      select
      </button>
      <button
        v-for="{ key, text } in shapes"
        :key="key"
        v-on:click="setShapeType(key)"
      >
        {{ text }}
      </button>
      <button v-on:click="removeAllShape()">清除</button>
    </div>
      <div id="draw-shape-k-line"
      class="k-line-chart"/>
</div>
</template>

<script>
import { init } from 'klinecharts'
import axios from 'axios'
const data = require('./demo.json')
const L1_COLOR = '#2196F3'
const L2_COLOR = '#FF9600'

export default {
  data () {
    return {
      code: data.code,
      level: 'day',
      shapes: [
        { key: 'priceLine', text: '价格线' },
        { key: 'priceChannelLine', text: '价格通道线' },
        { key: 'parallelStraightLine', text: '平行直线' },
        { key: 'fibonacciLine', text: '斐波那契回调' }
      ]
    }
  },
  mounted () {
    console.log()
    this.kLineChart = init('draw-shape-k-line')
    this.kLineChart.applyNewData(data.souce_data)
    this.kLineChart.createTechnicalIndicator('VOL')
    this.kLineChart.createTechnicalIndicator('MACD')
    this.init()
  },
  methods: {
    init () {
      const l0Arr = data.line[0]
      const l1Arr = data.line[1]

      // draw L1 line
      this.drawLine(l0Arr, L1_COLOR)

      // draw L2 line
      this.drawLine(l1Arr, L2_COLOR)
    },
    select_code () {
      console.log(this.code)
      axios.get('http://localhost:8000/result?code=' + this.code + '&level=' + this.level)
    },
    setShapeType: function (shapeName) {
      this.kLineChart.createShape(shapeName)
    },
    removeAllShape () {
      this.kLineChart.removeShape()
    },
    drawLine (arr, color) {
      arr.forEach(l => {
        let [y1, y2] = l.line_range
        const [x1, x2] = l.time_rangee

        const lineStyle = l.status === 'Grow' ? 'dash' : 'soild'
        if (l.direction === 'Down') {
          const tmp = y1
          y1 = y2
          y2 = tmp
        }
        this.createLine([
          { timestamp: x1, value: y1 },
          { timestamp: x2, value: y2 }
        ], color, lineStyle)
      })
    },
    createLine (points, color, style) {
      this.kLineChart.createShape(
        {
          name: 'segment',
          points,
          lock: true,
          mode: 'normal',
          styles: {
            line: {
              color,
              size: 2,
              style
            }
          }
        }
      )
    }
  }
}
</script>

<style>

</style>
