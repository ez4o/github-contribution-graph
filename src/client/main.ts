import * as d3 from "d3";
import { ContributionEntry } from './model/ContributionEntry';
// import { mockContributionData } from './data/mock_contribution_data';

// @ts-ignore
const hydratedContributionData = contributionData as unknown as ContributionEntry[];
// const hydratedContributionData: ContributionEntry[] = mockContributionData;

// @ts-ignore
const hydratedUsername = username;
// const hydratedUsername = 'Xyphuz';

// @ts-ignore
const hydratedImgBase64String = imgBase64String;

const width = 640
const height = 640

const margin = 35;
const textMargin = 10;
const chartMargin = 25;

const barOffsetX = 20;
const barWidth = 28;
const baseBarHeight = 10;

const fontSize = '0.7em';
const textWidth = 20;
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
  .range([width / 2 - margin - barOffsetX + chartMargin, width - margin * 2 - barOffsetX]);

const y = d3
  .scaleLinear()
  .range([height / 2, 0]);

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
  .selectAll('upper-bar')
  .data(hydratedContributionData)
  .enter()
  .append('rect')
  .attr('class', 'bar')
  .attr("clip-path", "url(#chart-background)")
  .attr('x', d => x(d.dateString)!)
  .attr('y', d => y(d.amount) - baseBarHeight / 2)
  .attr('width', barWidth)
  .attr('height', d => height / 2 - y(d.amount) + baseBarHeight)
  .attr('transform', `translate(${(x.bandwidth() - barWidth) / 2}, 0)`);

clipPath
  .selectAll('lower-bar')
  .data(hydratedContributionData)
  .enter()
  .append('rect')
  .attr('class', 'bar')
  .attr('x', d => x(d.dateString)!)
  .attr('y', height / 2 + baseBarHeight / 2)
  .attr('width', barWidth)
  .attr('height', d => height / 2 - y(d.amount) + baseBarHeight)
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
  .attr('y', d => y(d.amount) - lineHeight / 2)
  .attr('width', barWidth)
  .attr('height', d => height / 2 - y(d.amount))
  .attr('fill', '#000')
  .attr('transform', `translate(${(x.bandwidth() - barWidth) / 2}, 0)`)
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(d => d.amount);

svg
  .append('text')
  .attr('x', width / 2 - margin - barOffsetX + chartMargin - textMargin - textWidth / 2 - 4)
  .attr('y', height / 2 - lineHeight / 2)
  .attr('fill', '#000')
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(hydratedContributionData[0].dateString.substring(0, 4));

svg
  .append('text')
  .attr('x', width / 2 - margin - barOffsetX + chartMargin - textMargin - textWidth / 2 - 4)
  .attr('y', height / 2 + lineHeight / 2)
  .attr('fill', '#000')
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(hydratedContributionData[0].dateString.substring(5, 10));

svg
  .append('text')
  .attr('x', width - margin * 2 - barOffsetX + textMargin + textWidth / 2)
  .attr('y', height / 2 - lineHeight / 2)
  .attr('fill', '#000')
  .attr('text-anchor', 'middle')
  .attr('font-size', fontSize)
  .text(hydratedContributionData[hydratedContributionData.length - 1].dateString.substring(0, 4));

svg
  .append('text')
  .attr('x', width - margin * 2 - barOffsetX + textMargin + textWidth / 2)
  .attr('y', height / 2 + lineHeight / 2)
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