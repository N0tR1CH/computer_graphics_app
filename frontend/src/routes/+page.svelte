<script lang="ts">
	import TopBar from '$lib/components/top_bar.svelte';
	import ToolBar from '$lib/components/tool_bar.svelte';
	import ToolBarButton from '$lib/components/tool_bar_button.svelte';
	import TriangleOutline from '$lib/components/outlines/triangle_outline.svelte';
	import RectangleOutline from '$lib/components/outlines/rectangle_outline.svelte';
	import CircleOutline from '$lib/components/outlines/circle_outline.svelte';
	import MoveOutline from '$lib/components/outlines/move_outline.svelte';
	import ResizeOutline from '$lib/components/outlines/resize_outline.svelte';
	import StraightLineOutline from '$lib/components/outlines/straight_line_outline.svelte';
	import PencilOutline from '$lib/components/outlines/pencil_outline.svelte';
	import SaveOutline from '$lib/components/outlines/save_outline.svelte';
	import TextOutline from '$lib/components/outlines/text_outline.svelte';
	import CurveOutline from '$lib/components/outlines/curve_outline.svelte';
	import Canvas from '$lib/components/canvas.svelte';
	import Rectangle from '$lib/components/shapes/rectangle.svelte';
	import Triangle from '$lib/components/shapes/triangle.svelte';
	import Ellipse from '$lib/components/shapes/ellipse.svelte';
	import StraightLine from '$lib/components/shapes/straight_line.svelte';
	import Text from '$lib/components/shapes/text.svelte';
	import type { Shape } from '../types/shape';
	import type { PossibleActions } from '../types/possible_actions';
	import ColorPickers from '$lib/components/color_picking/color_pickers.svelte';
	import ThirdDimensionCanvas from '$lib/components/third_dimension_canvas.svelte';
	import Image from '$lib/components/shapes/image.svelte';
	import { main } from '$lib/wailsjs/go/models';
	import Swal from 'sweetalert2';
	import { fade } from 'svelte/transition';

	type NetPBMimg = {
		resource: string;
		comments: string[];
		status: string;
		base64str: string;
	};

	import { UploadNetPbmImg } from '$lib/wailsjs/go/main/Worker';
	import { EventsOnce } from '$lib/wailsjs/runtime/runtime';
	import {
		HandleRgbPointWiseTransformations,
		HandleAlphaPointWiseTransformations,
		HandleToGrayPointWiseTransformations,
		HandleFilterApplying,
		HandleHistogram,
		HandleBinarizeManual,
		HandleBinarizePercentBlack,
		HandleBinarizeOtsu,
		HandleBinarizeNiblack,
		HandleBinarizeBernsen,
		HandleDilation,
		HandleErosion,
		HandleOpening,
		HandleClosing
	} from '$lib/wailsjs/go/main/App';
	import { HandleBinarizeMeanIterative } from '$lib/wailsjs/go/main/App';
	import { HandleHitOrMiss } from '$lib/wailsjs/go/main/App';
	import BezierCurve from '$lib/components/shapes/bezier_curve.svelte';
	import TimelineOutline from '$lib/components/outlines/timeline_outline.svelte';
	import QuadraticCurve from '$lib/components/shapes/quadratic_curve.svelte';
	import DotsOutline from '$lib/components/outlines/dots_outline.svelte';

	window.addEventListener('keydown', (event) => {
		switch (event.key) {
			case 'ArrowLeft':
				shapes[shapes.length - 1].x = shapes[shapes.length - 1].x - 5;
				console.log('Left pressed');
				break;
			case 'ArrowRight':
				console.log('Right pressed');
				shapes[shapes.length - 1].x = shapes[shapes.length - 1].x + 5;
				break;
			case 'ArrowUp':
				console.log('Up pressed');
				shapes[shapes.length - 1].y = shapes[shapes.length - 1].y - 5;
				break;
			case 'ArrowDown':
				console.log('Down pressed');
				shapes[shapes.length - 1].y = shapes[shapes.length - 1].y + 5;
				break;
		}
	});

	let netpbmImages: NetPBMimg[] = [];
	let activeAction: PossibleActions = 'Triangle';
	let text: string = '';
	let shapes: Shape[] = [];
	let sceneWidth: number = 350;
	let sceneHeight: number = 350;

	let fileFormats: main.ImageFormat[] = [
		main.ImageFormat.jpg,
		main.ImageFormat.pbmP1,
		main.ImageFormat.pbmP4,
		main.ImageFormat.pgmP2,
		main.ImageFormat.pgmP5,
		main.ImageFormat.ppmP3,
		main.ImageFormat.ppmP6
	];
	let comments: string[] = [];
	let currentCommentInput: string = '';
	let selectedFileFormat: main.ImageFormat = main.ImageFormat.jpg;
</script>

<TopBar>
	<p class="text-2xl" slot="title">Stuff ;)</p>
</TopBar>
<ToolBar>
	<ToolBarButton bind:activeAction action={'Triangle'}>
		<TriangleOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Rectangle'}>
		<RectangleOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Circle'}>
		<CircleOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Move'}>
		<MoveOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Resize'}>
		<ResizeOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'StraightLine'}>
		<StraightLineOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Pencil'}>
		<PencilOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Save'}>
		<SaveOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Text'}>
		<TextOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Bezier'}>
		<TimelineOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'QuadraticCurve'}>
		<CurveOutline {activeAction} />
	</ToolBarButton>
	<ToolBarButton bind:activeAction action={'Polygon'}>
		<DotsOutline {activeAction} />
	</ToolBarButton>
</ToolBar>

<p class="text-bold text-center text-2xl text-white">
	Now arrows keys can also move images and shapes!
</p>

{#if activeAction == 'Save'}
	<div class="mx-auto my-4 max-w-sm" transition:fade>
		<label for="format" class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">
			File format options
		</label>
		<select
			bind:value={selectedFileFormat}
			on:change={() => (comments = [])}
			id="format"
			class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
		>
			{#each fileFormats as fileFormat}
				<option value={fileFormat}>
					{fileFormat}
				</option>
			{/each}
		</select>

		{#if selectedFileFormat != main.ImageFormat.jpg}
			<div class="my-4" transition:fade>
				<label for="comment" class="mb-2 block text-sm font-medium text-gray-900 dark:text-white"
					>Komentarz</label
				>
				<input
					type="text"
					bind:value={currentCommentInput}
					id="comment"
					class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
				/>
				<button
					on:click={() => {
						if (currentCommentInput != '') {
							comments = [...comments, currentCommentInput];
							currentCommentInput = '';
						}
					}}
					type="button"
					class="my-2 mb-2 w-full rounded-lg bg-purple-700 px-5 py-2.5 text-sm font-medium text-white hover:bg-purple-800 focus:outline-none focus:ring-4 focus:ring-purple-300 dark:bg-purple-600 dark:hover:bg-purple-700 dark:focus:ring-purple-900"
					>Dodaj komentarz do zapisywanego pliku</button
				>
			</div>
			{#each comments as comment}
				<span class="text-white">#{comment},&nbsp;</span>
			{/each}
		{/if}
	</div>
{/if}

<button
	type="button"
	class="my-4 mb-2 me-2 w-full rounded-full bg-blue-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
	on:click={async () => {
		const uuid = await UploadNetPbmImg();
		if (uuid == '') {
			return;
		}
		netpbmImages = [
			...netpbmImages,
			{ resource: uuid, comments: [], status: 'queued', base64str: '' }
		];
		EventsOnce(uuid, (comments, status, base64str) => {
			if (comments == null) {
				comments = [];
			}
			netpbmImages[netpbmImages.length - 1].comments = comments;
			netpbmImages[netpbmImages.length - 1].status = status;
			netpbmImages[netpbmImages.length - 1].base64str = base64str;
		});
	}}>Upload Image</button
>

{#if shapes[shapes.length - 1] !== undefined && shapes[shapes.length - 1].baseUrlImage !== ''}
	<div transition:fade>
		<button
			type="button"
			class="my-4 mb-2 me-2 w-full rounded-full bg-blue-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
			on:click={async () => {
				const { value } = await Swal.fire({
					title: 'Point-wise transformation',
					html: `
                        <form class="max-w-sm mx-auto">
                          <label for="actions" class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-600">Choose grayness method</label>
                          <select id="actions" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
                            <option selected value="rgb">RGB</option>
                            <option value="alpha">ALPHA</option>
                            <option value="gray">MAKE IT GRAY</option>
                          </select>
                        </form>
                          `,
					focusConfirm: false,
					preConfirm: () => document.getElementById('actions').value
				});
				switch (value) {
					case 'rgb': {
						const { value } = await Swal.fire({
							title: 'Point-wise transformation',
							html: `
                                <form class="max-w-sm mx-auto">
                          <label for="actions" class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-600">Choose an action</label>
                                  <select id="operation-select" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
                                    <option selected value="addition">Addition</option>
                                    <option value="substraction">Substraction</option>
                                    <option value="multiplication">Multiplication</option>
                                    <option value="division">Divide</option>
                                  </select>


                                  <label for="colors" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">Red</label>
                                <input id="red" type="number" min="0" max="255" value="0" />

                                  <label for="colors" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">Green</label>
                        <input id="green" type="number" min="0" max="255" value="0" />

                                  <label for="colors" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">Blue</label>
                        <input id="blue" type="number" min="0" max="255" value="0" />

                                </form>
                                  `,
							focusConfirm: false,
							preConfirm: () => {
								return [
									document.getElementById('operation-select').value,
									document.getElementById('red').value,
									document.getElementById('green').value,
									document.getElementById('blue').value
								];
							}
						});
						// Validate numbers
						for (let i = 1; i < 3; i++) {
							if (Number(value[i]) < 0 || Number(value[i]) > 255) {
								Swal.fire('RGB Numbers must be between 0 and 255!');
								break;
							}
						}
						const baseUrlImage = await HandleRgbPointWiseTransformations(
							value,
							shapes[shapes.length - 1].baseUrlImage
						);
						if (baseUrlImage == '') {
							console.error('baseUrlImage is empty');
							return;
						}
						shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
						break;
					}
					case 'alpha': {
						const { value } = await Swal.fire({
							title: 'Point-wise transformation',
							html: `
                                <form class="max-w-sm mx-auto">
                                    <label for="alpha" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">Alpha - transparency level</label>
                                    <input id="alpha" type="number" min="0" max="255" value="0" />
                                </form>
                                  `,
							focusConfirm: false,
							preConfirm: () => document.getElementById('alpha').value
						});

						if (Number(value) < 0 || Number(value) > 255) {
							Swal.fire('Alpha Number must be between 0 and 255!');
						}
						const baseUrlImage = await HandleAlphaPointWiseTransformations(
							Number(value),
							shapes[shapes.length - 1].baseUrlImage
						);
						if (baseUrlImage == '') {
							console.error('baseUrlImage is empty');
							return;
						}
						shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
						break;
					}
					case 'gray': {
						const { value } = await Swal.fire({
							title: 'Grayness method',
							html: `
                                <form class="max-w-sm mx-auto">
                          <label for="actions" class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-600">Choose an action</label>
                                  <select id="operation-select" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
                                    <option selected value="average">Average</option>
                                    <option value="weights">Weights</option>
                                  </select>
                                </form>
                                  `,
							focusConfirm: false,
							preConfirm: () => document.getElementById('operation-select').value
						});

						const baseUrlImage = await HandleToGrayPointWiseTransformations(
							value,
							shapes[shapes.length - 1].baseUrlImage
						);
						if (baseUrlImage == '') {
							console.error('baseUrlImage is empty');
							return;
						}
						shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
						break;
					}
					default: {
						Swal.fire('Problem!');
						break;
					}
				}
			}}>Apply point-wise intensity transformations to lastly inserted image</button
		>
	</div>

	<button
		type="button"
		class="my-4 mb-2 me-2 w-full rounded-full bg-blue-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
		on:click={async () => {
			const baseUrlImage = await HandleFilterApplying(shapes[shapes.length - 1].baseUrlImage);
			if (baseUrlImage == '') {
				console.error('baseUrlImage is empty');
				return;
			}
			shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
		}}
	>
		Apply filter
	</button>

	<button
		type="button"
		class="my-4 mb-2 me-2 w-full rounded-full bg-blue-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
		on:click={async () => {
			const baseUrlImage = await HandleHistogram(shapes[shapes.length - 1].baseUrlImage);
			if (baseUrlImage == '') {
				console.error('baseUrlImage is empty');
				return;
			}
			shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
		}}
	>
		Histogram actions
	</button>

	<button
		type="button"
		class="my-4 mb-2 me-2 w-full rounded-full bg-blue-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
		on:click={async () => {
			const { value } = await Swal.fire({
				title: 'Binarisation',
				html: `
                    <form class="max-w-sm mx-auto">
              <label for="actions" class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-600">Choose an action</label>
                      <select id="operation-select" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
                        <option selected value="manual">Manual</option>
                        <option value="percentblack">Percent Black selection</option>
                        <option value="meaniterative">Mean iterative selection</option>
                        <option value="otsu">Otsu</option>
                        <option value="niblack">Niblack</option>
                        <option value="bernsen">Bernsen</option>
                      </select>
                    </form>
                      `,
				focusConfirm: false,
				preConfirm: () => document.getElementById('operation-select').value
			});

			switch (value) {
				case 'manual': {
					const { value } = await Swal.fire({
						title: 'Manual method',
						html: `
                        <form class="max-w-sm mx-auto">
                            <label for="colors" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">Threshold</label>
                            <input id="threshold" type="number" min="0" max="255" value="0" />
                        </form>
                        `,
						focusConfirm: false,
						preConfirm: () => document.getElementById('threshold').value
					});
					const baseUrlImage = await HandleBinarizeManual(
						shapes[shapes.length - 1].baseUrlImage,
						Number(value)
					);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
				}

				case 'percentblack': {
					const { value } = await Swal.fire({
						title: 'Percent black',
						html: `
                        <form class="max-w-sm mx-auto">
                            <label for="colors" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">
                                Percent
                            </label>
                            <input id="threshold" type="number" min="0" max="100" value="0" step="0.01" />
                        </form>
                        `,
						focusConfirm: false,
						preConfirm: () => document.getElementById('threshold').value
					});
					const baseUrlImage = await HandleBinarizePercentBlack(
						shapes[shapes.length - 1].baseUrlImage,
						Number(value)
					);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
				}
				case 'meaniterative': {
					const { value } = await Swal.fire({
						title: 'Iterative mean',
						html: `
                        <form class="max-w-sm mx-auto">
                            <label for="colors" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">
                                Iterations
                            </label>
                            <input id="iterations" type="number" min="0" max="100" value="0" step="1" />
                        </form>
                        `,
						focusConfirm: false,
						preConfirm: () => document.getElementById('iterations').value
					});

					if (Number(value) < 0 || Number(value) > 100) {
						Swal.fire('Maximum number of iterations num must be between 0 and 100');
						return;
					}

					const baseUrlImage = await HandleBinarizeMeanIterative(
						shapes[shapes.length - 1].baseUrlImage,
						Number(value)
					);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
				}
				case 'otsu': {
					if (Number(value) < 0 || Number(value) > 100) {
						Swal.fire('Maximum number of iterations num must be between 0 and 100');
						return;
					}

					const baseUrlImage = await HandleBinarizeOtsu(shapes[shapes.length - 1].baseUrlImage);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
				case 'niblack': {
					const { value } = await Swal.fire({
						title: 'Niblack',
						html: `
                        <form class="max-w-sm mx-auto">
                            <label for="window-size" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">
                                Window Size
                            </label>
                            <input id="window-size" type="number" min="3" value="0" step="1" />
                            <label for="k-val" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">
                                Value of k
                            </label>
                            <input id="k-val" type="number" step="1"/>
                        </form>
                        `,
						focusConfirm: false,
						preConfirm: () => [
							document.getElementById('window-size').value,
							document.getElementById('k-val').value
						]
					});
					const baseUrlImage = await HandleBinarizeNiblack(
						shapes[shapes.length - 1].baseUrlImage,
						Number(value[0]),
						Number(value[1])
					);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
				case 'bernsen': {
					const { value } = await Swal.fire({
						title: 'Grayness method',
						html: `
                        <form class="max-w-sm mx-auto">
                            <label for="window-size" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">
                                Window Size
                            </label>
                            <input id="window-size" type="number" min="3" value="0" step="1" />
                            <label for="contrast" class="block mb-2 text-sm font-medium text-gray-900 dark:text-black">
                                Contrast value
                            </label>
                            <input id="contrast" type="number" step="0.01"/>
                        </form>
                        `,
						focusConfirm: false,
						preConfirm: () => [
							document.getElementById('window-size').value,
							document.getElementById('contrast').value
						]
					});
					const baseUrlImage = await HandleBinarizeBernsen(
						shapes[shapes.length - 1].baseUrlImage,
						Number(value[0]),
						Number(value[1])
					);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
			}
		}}
	>
		Binarisation
	</button>

	<button
		type="button"
		class="my-4 mb-2 me-2 w-full rounded-full bg-blue-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
		on:click={async () => {
			const { value } = await Swal.fire({
				title: 'Morphological filters',
				html: `
                    <form class="max-w-sm mx-auto">
                        <label
                            for="morphological-filters"
                            class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-600"
                        >
                            Choose an action
                        </label>
                        <select
                            id="morphological-filters-select"
                            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                        >
                            <option selected value="dilation">
                                Dilation
                            </option>
                            <option value="erosion">
                                Erosion
                            </option>
                            <option value="opening">
                                Opening
                            </option>
                            <option value="closing">
                                Closing
                            </option>
                            <option value="hit-or-miss">
                                Hit or Miss
                            </option>
                        </select>
                    </form>
            `,
				focusConfirm: false,
				preConfirm: () => document.getElementById('morphological-filters-select').value
			});

			switch (value) {
				case 'dilation': {
					const baseUrlImage = await HandleDilation(shapes[shapes.length - 1].baseUrlImage);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
				case 'erosion': {
					const baseUrlImage = await HandleErosion(shapes[shapes.length - 1].baseUrlImage);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
				case 'opening': {
					const baseUrlImage = await HandleOpening(shapes[shapes.length - 1].baseUrlImage);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
				case 'closing': {
					const baseUrlImage = await HandleClosing(shapes[shapes.length - 1].baseUrlImage);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
				case 'hit-or-miss': {
					const baseUrlImage = await HandleHitOrMiss(shapes[shapes.length - 1].baseUrlImage);
					if (baseUrlImage == '') {
						console.error('baseUrlImage is empty');
						return;
					}
					shapes[shapes.length - 1].baseUrlImage = baseUrlImage;
					break;
				}
			}
		}}
	>
		Morphological filters
	</button>
{/if}

<div class="relative overflow-x-auto">
	<table class="my-4 w-full text-left text-sm text-gray-500 rtl:text-right dark:text-gray-400">
		<thead class="bg-gray-50 text-xs uppercase text-gray-700 dark:bg-gray-700 dark:text-gray-400">
			<tr>
				<th scope="col" class="px-6 py-3">Resource</th>
				<th scope="col" class="px-6 py-3">Comments</th>
				<th scope="col" class="px-6 py-3">Status</th>
				<th scope="col" class="px-6 py-3">Action</th>
			</tr>
		</thead>
		<tbody>
			{#each netpbmImages as netpbmImage}
				<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800" transition:fade>
					<th
						scope="row"
						class="whitespace-nowrap px-6 py-4 font-medium text-gray-900 dark:text-white"
					>
						{netpbmImage.resource}
					</th>
					<td class="px-6 py-4">
						{#each netpbmImage.comments as comment}
							{comment},&nbsp;
						{/each}
					</td>
					<td class="px-6 py-4">{netpbmImage.status}</td>
					<td class="px-6 py-4">
						{#if netpbmImage.status == 'completed'}
							<button
								on:click={() => {
									shapes = [
										...shapes,
										{
											name: 'Image',
											x: 100,
											y: 100,
											height: 0,
											width: 0,
											base: 0,
											radius1: 0,
											radius2: 0,
											rotation: 0,
											x1: 0,
											y1: 0,
											text: '',
											hexColor: '',
											baseUrlImage: `data:image/png;base64,${netpbmImage.base64str}`
										}
									];
								}}
								type="button"
								class="mb-2 me-2 rounded-full bg-blue-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
								>Add</button
							>
						{/if}
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<Canvas
	height={500}
	width={1240}
	bind:shapes
	bind:activeAction
	bind:text
	bind:selectedFileFormat
	bind:comments
>
	{#each shapes as shape}
		{#if shape.name === 'Rectangle'}
			<Rectangle
				x={shape.x}
				y={shape.y}
				height={shape.height}
				width={shape.width}
				hexColor={shape.hexColor}
			/>
		{:else if shape.name === 'Triangle'}
			<Triangle
				x={shape.x}
				y={shape.y}
				base={shape.base}
				height={shape.height}
				hexColor={shape.hexColor}
			/>
		{:else if shape.name === 'Ellipse'}
			<Ellipse
				x={shape.x}
				y={shape.y}
				radius1={shape.radius1}
				radius2={shape.radius2}
				rotation={shape.rotation}
				hexColor={shape.hexColor}
			/>
		{:else if shape.name === 'StraightLine'}
			<StraightLine x={shape.x} y={shape.y} x1={shape.x1} y1={shape.y1} hexColor={shape.hexColor} />
		{:else if shape.name === 'Text'}
			<Text x={shape.x} y={shape.y} text={shape.text} hexColor={shape.hexColor} />
		{:else if shape.name === 'Image'}
			<Image x={shape.x} y={shape.y} baseUrlImage={shape.baseUrlImage} />
		{:else if shape.name === 'Bezier'}
			<BezierCurve
				hexColor={shape.hexColor}
				start={shape.bezierStart}
				cp1={shape.bezierCp1}
				cp2={shape.bezierCp2}
				end={shape.bezierEnd}
			/>
		{:else if shape.name === 'QuadraticCurve'}
			<QuadraticCurve
				hexColor={shape.hexColor}
				start={shape.bezierStart}
				cp={shape.bezierCp1}
				end={shape.bezierEnd}
			/>
		{/if}
	{/each}
</Canvas>

{#if activeAction === 'Move'}
	<p class="mt-4 text-center text-white" transition:fade>
		Current mode is moving elements. Press left click in order to move latest element somewhere else
		on the canvas.
	</p>
{:else if activeAction === 'Resize'}
	<p class="mt-4 text-center text-white" transition:fade>
		Current mode is resizing elements. Press left click in order to scale latest element. If it
		doesn't do anything that means element is not scalable!
	</p>
{:else if activeAction === 'Save'}
	<p class="mt-4 text-center text-white" transition:fade>
		Current mode is saving canvas. Click on the canvas in order to invoke action!
	</p>
{:else if activeAction === 'Text'}
	<p class="mt-4 text-center text-white" transition:fade>
		Write text to input field and click on the canvas in order to add it to correct coordinates.
	</p>
	<input
		class="block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500"
		bind:value={text}
		type="text"
	/>
{:else if (activeAction === 'Bezier' && shapes[shapes.length - 1]?.name === 'Bezier') || (activeAction === 'QuadraticCurve' && shapes[shapes.length - 1]?.name === 'QuadraticCurve')}
	<div>
		<div class="my-4">
			<label for="small-input" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>START</label
			>
			<div class="flex gap-x-2">
				<label for="" class="text-white">X</label>
				<input
					bind:value={shapes[shapes.length - 1].bezierStart.x}
					type="number"
					id="small-input"
					class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
				/>
				<label for="" class="text-white">Y</label>
				<input
					bind:value={shapes[shapes.length - 1].bezierStart.y}
					type="number"
					id="small-input"
					class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
				/>
			</div>
		</div>

		<div>
			<label for="small-input" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>END</label
			>
			<div class="flex gap-x-2">
				<label for="" class="text-white">X</label>
				<input
					bind:value={shapes[shapes.length - 1].bezierEnd.x}
					type="number"
					id="small-input"
					class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
				/>
				<label for="" class="text-white">Y</label>
				<input
					bind:value={shapes[shapes.length - 1].bezierEnd.y}
					type="number"
					id="small-input"
					class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
				/>
			</div>
		</div>

		<div>
			<label for="small-input" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
				>CONTROL POINT 1</label
			>
			<div class="flex gap-x-2">
				<label for="" class="text-white">X</label>
				<input
					bind:value={shapes[shapes.length - 1].bezierCp1.x}
					type="number"
					id="small-input"
					class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
				/>
				<label for="" class="text-white">Y</label>
				<input
					bind:value={shapes[shapes.length - 1].bezierCp1.y}
					type="number"
					id="small-input"
					class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
				/>
			</div>
		</div>
		{#if shapes[shapes.length - 1].name === 'Bezier'}
			<div>
				<label
					for="small-input"
					class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
					>CONTROL POINT 2</label
				>
				<div class="flex gap-x-2">
					<label for="" class="text-white">X</label>
					<input
						bind:value={shapes[shapes.length - 1].bezierCp2.x}
						type="number"
						id="small-input"
						class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
					/>
					<label for="" class="text-white">Y</label>
					<input
						bind:value={shapes[shapes.length - 1].bezierCp2.y}
						type="number"
						id="small-input"
						class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-xs focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
					/>
				</div>
			</div>
		{/if}
	</div>
{/if}

<ColorPickers></ColorPickers>
{#key sceneWidth + sceneHeight}
	<ThirdDimensionCanvas bind:shapes width={sceneWidth} height={sceneHeight}></ThirdDimensionCanvas>
{/key}
<div class="mb-4 flex justify-center gap-x-4">
	<label for="" class="font-bold text-white">
		Width of the scene
		<input type="number" bind:value={sceneWidth} class="rounded-xl p-2 text-black" />
	</label>
	<label for="" class="font-bold text-white">
		Height of the scene
		<input type="number" bind:value={sceneHeight} class="rounded-xl p-2 text-black" />
	</label>
</div>
<p class="mb-4 text-center font-bold text-white">
	Cube is being added to fixed position but you can change it by using move tool although it is
	clunky (3d to 2d ehh)
</p>

<style lang="postcss">
	:global(html) {
		background-color: rgba(0, 0, 0, 0);
	}

	select {
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
	}

	select::-ms-expand {
		display: none;
	}
</style>
