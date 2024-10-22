<script lang="ts">
	import iro from '@jaames/iro';
	import { onMount } from 'svelte';
	import { RgbToCmyk, CmykToRgb } from '$lib/wailsjs/go/main/App';
	import { currentColor } from '$lib/stores/stores';

	let settingRgb: boolean = false;
	let settingCmyk: boolean = false;
	let colorPicker: iro.ColorPicker;
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

		$currentColor = colorPicker.color.hexString;
		r = colorPicker.color.red;
		g = colorPicker.color.green;
		b = colorPicker.color.blue;
		h = colorPicker.color.hue;
		s = colorPicker.color.saturation;
		v = colorPicker.color.value;
		setCmyk();

		colorPicker.on('color:change', () => {
			$currentColor = colorPicker.color.hexString;
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

<h2 class="mt-2 text-center text-4xl font-bold text-white">RGB/HSV</h2>
<div class="mt-4 flex justify-around gap-x-4">
	<div id="picker"></div>
	<ul>
		<li class="mb-4 font-bold text-white">
			<input type="number" min="0" max="255" class="rounded-full px-2 text-black" bind:value={r} />
		</li>
		<li class="mb-4 font-bold text-white">
			<input type="number" min="0" max="255" class="rounded-full px-2 text-black" bind:value={g} />
		</li>
		<li class="mb-4 font-bold text-white">
			<input type="number" min="0" max="255" class="rounded-full px-2 text-black" bind:value={b} />
		</li>
		<li class="mb-4 font-bold text-white">
			<input type="number" min="0" max="360" class="rounded-full px-2 text-black" bind:value={h} />
		</li>
		<li class="mb-4 font-bold text-white">
			<input type="number" min="0" max="100" class="rounded-full px-2 text-black" bind:value={s} />
		</li>
		<li class="mb-4 font-bold text-white">
			<input type="number" min="0" max="100" class="rounded-full px-2 text-black" bind:value={v} />
		</li>
	</ul>
	<ul>
		<li class="mb-4 font-bold text-white">RED</li>
		<li class="mb-4 font-bold text-white">GREEN</li>
		<li class="mb-4 font-bold text-white">BLUE</li>
		<li class="mb-4 font-bold text-white">HUE</li>
		<li class="mb-4 font-bold text-white">SATURATION</li>
		<li class="mb-4 font-bold text-white">VALUE</li>
	</ul>
</div>

<div>
	<h2 class="text-center text-4xl font-bold text-white">CMYK</h2>
	<ul class="flex justify-around">
		<li>
			<span class="font-bold text-white">C</span>
			<input type="number" min="0" max="255" class="rounded-full px-2 text-black" bind:value={c} />
		</li>
		<li>
			<span class="font-b255 text-white">M</span>
			<input type="number" min="0" max="255" class="rounded-full px-2 text-black" bind:value={m} />
		</li>
		<li>
			<span class="font-b255 text-white">Y</span>
			<input type="number" min="0" max="255" class="rounded-full px-2 text-black" bind:value={y} />
		</li>
		<li>
			<span class="font-b255 text-white">K</span>
			<input type="number" min="0" max="255" class="rounded-full px-2 text-black" bind:value={k} />
		</li>
	</ul>
</div>

<h2 class="mt-2 text-center text-4xl font-bold text-white">Color Preview</h2>
<div
	style="background-color: {$currentColor}"
	class="my-2 h-12 rounded-full border-2 border-solid border-white"
></div>
