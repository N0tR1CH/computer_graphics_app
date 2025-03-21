<script lang="ts">
	import { onMount } from 'svelte';
	import * as THREE from 'three';
	import { OrbitControls } from 'three/addons/controls/OrbitControls.js';
	import { currentColor } from '$lib/stores/stores';
	import type { Shape } from '../../types/shape';

	export let shapes: Shape[];
	export let width: number = 350;
	export let height: number = 350;
	let canvasContainer: HTMLElement;
	let renderer: THREE.WebGLRenderer;
	let controls: OrbitControls;

	const DownloadCanvasAsImage = () => {
		const baseUrlImage = renderer.domElement.toDataURL('image/png');
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
				hexColor: $currentColor,
				baseUrlImage: baseUrlImage
			}
		];
	};

	onMount(() => {
		// Huge credits: https://codepen.io/aexklon/pen/mdGmmxJ
		// Width and Height of the threejs canvas
		const [w, h] = [width, height];
		// WebGL Rendering Engine
		renderer = new THREE.WebGLRenderer({
			alpha: true,
			antialias: true,
			preserveDrawingBuffer: true
		});

		renderer.domElement.classList.add(
			'border-solid',
			'border-2',
			'border-white',
			'rounded-xl',
			'm-4'
		);

		renderer.setSize(w, h);
		renderer.shadowMap.enabled = true;
		canvasContainer.appendChild(renderer.domElement);

		// New Scene
		const scene = new THREE.Scene();

		// Perspective Camera
		const camera = new THREE.PerspectiveCamera(45, w / h, 0.1, 150);

		// Position x,y,z axis of camera
		camera.position.set(0, 0, 3);

		// Controls
		controls = new OrbitControls(camera, renderer.domElement);

		// create the mesh from geometry and material
		const geometry = new THREE.BoxGeometry(1, 1, 1);

		// Create an array to store colors for each vertex
		const colors = [];
		const positionAttribute = geometry.getAttribute('position');

		// Assign random colors based on vertex positions
		for (let i = 0; i < positionAttribute.count; i++) {
			const x = positionAttribute.getX(i);
			const y = positionAttribute.getY(i);
			const z = positionAttribute.getZ(i);

			// Generate color based on position (RGB)
			const color = new THREE.Color(x + 0.5, y + 0.5, z + 0.5);
			colors.push(color.r, color.g, color.b);
		}

		// Set colors
		geometry.setAttribute('color', new THREE.Float32BufferAttribute(colors, 3));
		const material = new THREE.MeshBasicMaterial({ vertexColors: true });
		const cube = new THREE.Mesh(geometry, material);
		cube.rotation.x = -Math.PI / 4;
		cube.rotation.z = Math.PI / 4;
		scene.add(cube);

		// lights
		const light = new THREE.DirectionalLight(0xeaeaea);
		light.position.set(5, 25, 50);
		light.castShadow = true;

		// adding objects to the scene
		scene.add(cube);
		scene.add(light);

		// render the scene
		const render = () => {
			requestAnimationFrame(render);

			controls.update();

			renderer.render(scene, camera);
		};

		render();

		return () => {
			renderer.dispose();
			controls.dispose();
		};
	});
</script>

<div
	class="flex justify-center items-center border-solid border-white border-2 rounded-xl mb-2"
	bind:this={canvasContainer}
>
	<button
		type="button"
		class="text-white bg-blue-700 hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 font-medium rounded-full text-sm px-5 py-2.5 text-center me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
		on:click={() => {
			DownloadCanvasAsImage();
		}}>Add to canvas</button
	>
</div>
