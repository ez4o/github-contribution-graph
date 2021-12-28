import './style.css'
import * as d3 from "d3";
import { mockContributionData } from './data/mock_contribution_data';

const width = window.innerWidth / 2
const height = window.innerHeight / 2

const svg = d3.select('#app')
  .append('svg')
  .attr('width', width)
  .attr('height', height);

const x = d3
  .scaleBand()
  .range([0, width])

const y = d3.scaleLinear().range([height, 0])

function getHeightEndpoint(val: number) {
  let count = Math.floor(val).toString().length - 1
  let step = Math.pow(10, count)

  if (val / step < 5) {
    step = step / 2
  }

  count = Math.ceil(val / step)

  return count * step
}

const endPoint = getHeightEndpoint(d3.max(mockContributionData.map(d => d.amount))!)

x.domain(mockContributionData.map(d => d.date.toISOString()))
y.domain([0, endPoint])

svg
  .selectAll('rect')
  .data(mockContributionData)
  .enter()
  .append('rect')
  .attr('x', d => x(d.date.toISOString())!)
  .attr('width', 24)
  .attr('y', d => y(d.amount))
  .attr('height', d => height - y(d.amount))
  .attr('fill', '#09c')
  .attr('transform', `translate(${x.bandwidth() / 2 - 12}, 0)`)