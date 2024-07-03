<script>
    import { onMount } from "svelte";
    // Pepe is scared of the matrix
    import { userInfo } from "../stores";
    import { get } from "svelte/store"
    //
    
    let canvas;

    onMount(() => {
        // pepe & cat is scared of the matrix
        var email = get(userInfo).email
        if (email === "pepe" || email === "cat") {
            return
        }
        //
        var c = canvas;
        var ctx = c.getContext("2d");

        // Making the canvas full screen
        c.height = window.innerHeight;
        c.width = window.innerWidth;

        // Characters for the matrix effect
        var matrix = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789@#$%^&*()*&^%+-/~{[|`]}";
        matrix = matrix.split("");

        var font_size = 14;
        var columns = c.width / font_size;
        var drops = [];

        for (var x = 0; x < columns; x++)
            drops[x] = 1;

        function draw() {
            ctx.fillStyle = "rgba(0,0,0, 0.04)";
            ctx.fillRect(0, 0, c.width, c.height);

            ctx.fillStyle = "greenyellow";
            ctx.font = font_size + "px arial";

            for (var i = 0; i < drops.length; i++) {
                var text = matrix[Math.floor(Math.random() * matrix.length)];
                ctx.fillText(text, i * font_size, drops[i] * font_size);

                if (drops[i] * font_size > c.height && Math.random() > 0.975)
                    drops[i] = 0;

                drops[i]++;
            }
        }

        setInterval(draw, 45);
    });
</script>
<main>
    <canvas bind:this={canvas}></canvas>
</main>

<style>
    
    canvas {
        display: block;
        width: 100%;
        height: 100%;
    }
</style>
