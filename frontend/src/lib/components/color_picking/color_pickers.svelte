<script lang="ts">
	import iro from '@jaames/iro';
	import { onMount } from 'svelte';
	import { RgbToCmyk, CmykToRgb } from '$lib/wailsjs/go/main/App';

	let settingRgb: boolean = false;
	let settingCmyk: boolean = false;
	let colorPicker: iro.ColorPicker;
	let hexColor: string;
	let r: number = 0;
	let g: number = 0;
	let b: number = 0;
	let h: number = 0;
	let s: number = 0;
	let v: number = 0;
	let c: number = 0;
	let m: number = 0;
	let y: number = 0;
	let k: number = 0;

	async function setCmyk() {
		if (settingCmyk) {
			return;
		}

		settingCmyk = true;
		const cmyk = await RgbToCmyk(r, g, b);
		settingCmyk = false;
		c = cmyk.c;
		m = cmyk.m;
		y = cmyk.y;
		k = cmyk.k;
	}

	async function setRgb() {
		if (settingRgb) {
			return;
		}

		settingRgb = true;
		const rgb = await CmykToRgb(c, m, y, k);
		settingRgb = false;
		colorPicker.color.red = rgb.r;
		colorPicker.color.green = rgb.g;
		colorPicker.color.blue = rgb.b;
	}

	$: if (r) {
		colorPicker.color.red = r;
	}

	$: if (g) {
		colorPicker.color.green = g;
	}

	$: if (b) {
		colorPicker.color.blue = b;
	}

	$: if (h) {
		colorPicker.color.hue = h;
	}

	$: if (s) {
		colorPicker.color.saturation = s;
	}

	$: if (v) {
		colorPicker.color.value = v;
	}

	$: if (c || m || y || k) {
		setRgb();
	}

	onMount(() => {
		colorPicker = new iro.ColorPicker('#picker', {
			layout: [
				{
					component: iro.ui.Slider,
					options: {
						sliderType: 'red',
						borderWidth: 1
					}
				},
				{
					component: iro.ui.Slider,
					options: {
						sliderType: 'green',
						borderWidth: 1
					}
				},
				{
					component: iro.ui.Slider,
					options: {
						sliderType: 'blue',
						borderWidth: 1
					}
				},
				{
					component: iro.ui.Slider,
					options: {
						sliderType: 'hue',
						borderWidth: 1
					}
				},
				{
					component: iro.ui.Slider,
					options: {
						sliderType: 'saturation',
						borderWidth: 1
					}
				},
				{
					component: iro.ui.Slider,
					options: {
						sliderType: 'value',
						borderWidth: 1
					}
				}
			]
		});

		hexColor = colorPicker.color.hexString;
		r = colorPicker.color.red;
		g = colorPicker.color.green;
		b = colorPicker.color.blue;
		h = colorPicker.color.hue;
		s = colorPicker.color.saturation;
		v = colorPicker.color.value;
		setCmyk();

		colorPicker.on('color:change', () => {
			hexColor = colorPicker.color.hexString;
			r = colorPicker.color.red;
			g = colorPicker.color.green;
			b = colorPicker.color.blue;
			h = colorPicker.color.hue;
			s = colorPicker.color.saturation;
			v = colorPicker.color.value;
			setCmyk();
		});

		return () => {
			colorPicker.off('color:change', () => {});
		};
	});
</script>

<h2 class="text-center text-4xl text-white font-bold mt-2">RGB/HSV</h2>
<div class="flex gap-x-4 justify-around mt-4">
	<div id="picker"></div>
	<ul>
		<li class="text-white mb-4 font-bold">
			<input type="number" min="0" max="255" class="rounded-full text-black px-2" bind:value={r} />
		</li>
		<li class="text-white mb-4 font-bold">
			<input type="number" min="0" max="255" class="rounded-full text-black px-2" bind:value={g} />
		</li>
		<li class="text-white mb-4 font-bold">
			<input type="number" min="0" max="255" class="rounded-full text-black px-2" bind:value={b} />
		</li>
		<li class="text-white mb-4 font-bold">
			<input type="number" min="0" max="360" class="rounded-full text-black px-2" bind:value={h} />
		</li>
		<li class="text-white mb-4 font-bold">
			<input type="number" min="0" max="100" class="rounded-full text-black px-2" bind:value={s} />
		</li>
		<li class="text-white mb-4 font-bold">
			<input type="number" min="0" max="100" class="rounded-full text-black px-2" bind:value={v} />
		</li>
	</ul>
	<ul>
		<li class="text-white mb-4 font-bold">RED</li>
		<li class="text-white mb-4 font-bold">GREEN</li>
		<li class="text-white mb-4 font-bold">BLUE</li>
		<li class="text-white mb-4 font-bold">HUE</li>
		<li class="text-white mb-4 font-bold">SATURATION</li>
		<li class="text-white mb-4 font-bold">VALUE</li>
	</ul>
</div>

<div>
	<h2 class="text-center text-4xl text-white font-bold">CMYK</h2>
	<ul class="flex justify-around">
		<li>
			<span class="text-white font-bold">C</span>
			<input type="number" min="0" max="255" class="rounded-full text-black px-2" bind:value={c} />
		</li>
		<li>
			<span class="text-white font-b255">M</span>
			<input type="number" min="0" max="255" class="rounded-full text-black px-2" bind:value={m} />
		</li>
		<li>
			<span class="text-white font-b255">Y</span>
			<input type="number" min="0" max="255" class="rounded-full text-black px-2" bind:value={y} />
		</li>
		<li>
			<span class="text-white font-b255">K</span>
			<input type="number" min="0" max="255" class="rounded-full text-black px-2" bind:value={k} />
		</li>
	</ul>
</div>

<h2 class="text-center text-4xl text-white font-bold mt-2">Color Preview</h2>
<div
	style="background-color: {hexColor}"
	class="h-12 my-2 rounded-full border-white border-solid border-2"
></div>
