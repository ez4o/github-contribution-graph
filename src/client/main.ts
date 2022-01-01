import * as d3 from "d3";
import { ContributionEntry } from './model/contribution_entry';
import { mockContributionData } from './data/mock_contribution_data';
import { mockImgBase64String } from "./data/mock_img_base64_string";

const DEBUG = false;

let hydratedContributionData: ContributionEntry[];
let hydratedUsername: string;
let hydratedImgBase64String: string;

if (DEBUG) {
  hydratedContributionData = mockContributionData;
  hydratedUsername = 'Xyphuz';
  hydratedImgBase64String = mockImgBase64String;
} else {
  // @ts-ignore
  hydratedContributionData = contributionData as unknown as ContributionEntry[];
  // @ts-ignore
  hydratedUsername = username;
  // @ts-ignore
  hydratedImgBase64String = imgBase64String;
}

const width = 640
const height = 640

const margin = 35;
const chartScaleMarginX = 35;
const chartScaleMarginY = 0;

const textMargin = 10;
const startDateTextOffsetX = 31;

const barOffsetX = 20;
const barWidth = 28;
const baseBarHeight = 10;

const fontSize = '0.7em';
const textWidth = 20;
const baseOffsetY = 8;
const lineHeight = 20;
const titleLineHeight = 40;

const titleX = 80;
const titleY = 500;

const svg = d3
  .select('#svg-container')
  .append('svg')
  .attr('width', width)
  .attr('height', height);

svg
  .append('svg:image')
  .attr('id', 'background')
  .attr('width', width)
  .attr('height', height)
  .attr('xlink:href', hydratedImgBase64String)
  .attr('preserveAspectRatio', 'none');

const x = d3
  .scaleBand()
  .range([width / 2 - margin - barOffsetX + chartScaleMarginX, width - margin - barOffsetX - chartScaleMarginX]);

const y = d3
  .scaleLinear()
  .range([margin + chartScaleMarginY, height - margin - chartScaleMarginY]);

const endPoint = ((max: number) => {
  return (Math.ceil(max / 4) + 1) * 4;
})(d3.max(hydratedContributionData.map(d => d.amount))!);

x.domain(hydratedContributionData.map(d => d.dateString));
y.domain([0, endPoint]);

svg
  .append('rect')
  .attr('id', 'chart-background')
  .attr('width', width - margin * 2)
  .attr('height', height - margin * 2)
  .attr('fill', '#fff')
  .attr('x', margin)
  .attr('y', margin);

const clipPath = svg
  .append("defs")
  .append("clipPath")
  .attr("id", "clip")

clipPath
  .selectAll('bar')
  .data(hydratedContributionData)
  .enter()
  .append('rect')
  .attr('class', 'bar')
  .attr('x', d => x(d.dateString)!)
  .attr('y', d => height / 2 - (y(d.amount) + baseBarHeight) / 2)
  .attr('width', barWidth)
  .attr('height', d => y(d.amount) + baseBarHeight)
  .attr('fill', `rgba(0, 0, 0, 0.2)`)
  .attr('transform', `translate(${(x.bandwidth() - barWidth) / 2}, 0)`);

svg
  .append('svg:image')
  .attr("clip-path", "url(#clip)")
  .attr('width', width)
  .attr('height', height)
  .attr('xlink:href', hydratedImgBase64String)
  .attr('preserveAspectRatio', 'none');

svg
  .selectAll('contribution-amount')
  .data(hydratedContributionData)
  .enter()
  .append('text')
  .attr('x', d => x(d.dateString)! + barWidth / 2)
  .attr('y', d => height / 2 - (y(d.amount) + baseBarHeight) / 2 - lineHeight / 2)
  .attr('width', barWidth)
  .attr('height', d => height / 2 - y(d.amount))
  .attr('fill', '#000')
  .attr('transform', `translate(${(x.bandwidth() - barWidth) / 2}, 0)`)
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(d => d.amount);

svg
  .append('text')
  .attr('x', width / 2 - margin - barOffsetX + startDateTextOffsetX - textMargin - textWidth / 2)
  .attr('y', height / 2 - lineHeight / 2 + baseOffsetY)
  .attr('fill', '#000')
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(hydratedContributionData[0].dateString.substring(0, 4));

svg
  .append('text')
  .attr('x', width / 2 - margin - barOffsetX + startDateTextOffsetX - textMargin - textWidth / 2)
  .attr('y', height / 2 + lineHeight / 2 + baseOffsetY)
  .attr('fill', '#000')
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(hydratedContributionData[0].dateString.substring(5, 10));

svg
  .append('text')
  .attr('x', width - margin * 2 - barOffsetX + textMargin + textWidth / 2)
  .attr('y', height / 2 - lineHeight / 2 + baseOffsetY)
  .attr('fill', '#000')
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(hydratedContributionData[hydratedContributionData.length - 1].dateString.substring(0, 4));

svg
  .append('text')
  .attr('x', width - margin * 2 - barOffsetX + textMargin + textWidth / 2)
  .attr('y', height / 2 + lineHeight / 2 + baseOffsetY)
  .attr('fill', '#000')
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(hydratedContributionData[hydratedContributionData.length - 1].dateString.substring(5, 10));

svg
  .append('text')
  .attr('x', titleX)
  .attr('y', titleY)
  .attr('fill', '#666')
  .attr('font-size', fontSize)
  .text('Contribution Graph of');

svg
  .append('text')
  .attr('x', titleX)
  .attr('y', titleY + titleLineHeight)
  .attr('fill', '#000')
  .attr('font-size', "2rem")
  .attr('font-weight', 'bold')
  .text(hydratedUsername);