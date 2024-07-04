<script>
    export let styleConfig;
    export let onClick;
    export let btnText;
    function parseStyle(config) {
        return `
            --btn-type: ${config.btnType};
            --btn-width: ${config.btnWidth}px;
            --btn-height: ${config.btnHeight}px;
            --font-size: ${config.fontSize}px;
            --font-weight: ${config.fontWeight};
            --font-color: ${config.fontColor};
            --font-hover-color: ${config.hoverFontColor};
            --border-radius: ${config.borderRadius}px;
            --border-width: ${config.borderWidth}px;
            --border-color: ${config.borderColor};
            --border-hover-color: ${config.borderHoverColor};
            --background1-color: ${config.backgroundTone1};
            --background2-color: ${config.backgroundTone2};
            --background3-color: ${config.backgroundTone3};
        `;
    }

    const holdDuration = 1200;
    let timer;
    let progress = 0;
    let isHolding = false;

    function startTimer() {
        isHolding = true;
        progress = 0;
        timer = setInterval(() => {
            progress += 10; // Every 10ms
            if (progress >= holdDuration) {
                clearInterval(timer);
                progress = holdDuration;
                isHolding = false;
                onClick();
            }
        }, 10);
    }

    function stopTimer() {
        clearInterval(timer);
        progress = 0;
        isHolding = false;
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div 
    class="btn btn-following"
    style={parseStyle(styleConfig)}
    on:mousedown={styleConfig.btnType === 2 ? startTimer : null}
    on:mouseup={stopTimer}
    on:mouseleave={stopTimer}
    on:click={styleConfig.btnType !== 2 ? onClick : null}
>
    <span></span><span></span><span></span><span></span>
    <span></span><span></span><span></span><span></span>
    <span></span><span></span><span></span>
    {#if styleConfig.btnType === 2}
        <div class="progress" style="width: {progress / holdDuration * 100}%"></div>
    {/if}
    <h4>{btnText}</h4>
</div>

<style>
    * {
        padding: 0;
        margin: 0;
        box-sizing: border-box;
    }

    h4 {
        text-wrap: nowrap;
        user-select: none;
        font-weight: var(--font-weight, 600);
        font-size: var(--font-size, 20px);
        color: var(--font-color, red);
        transition: color ease-in-out 0.5s;
    }

    .btn {
        cursor: pointer;
        height: var(--btn-height, 60px);
        width: var(--btn-width, 150px);
        display: flex;
        justify-content: center;
        align-items: center;
        transition: all 0.8s ease;
        position: relative;
        border: solid #111;
        border-width: var(--border-width, 2px);
        border-color: var(--border-color, white);
        border-radius: var(--border-radius, 0px);
        background-color: var(--background1-color, grey);
        transition: border-color ease-in-out 0.5s;
    }

    .btn::before,
    .btn::after {
        position: absolute;
        display: block;
        content: "";
        width: 100%;
        height: 100%;
    }

    .progress {
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        background-color: rgba(155, 138, 221, 0.714);
        transition: width 0.01s linear;
        pointer-events: none;
        z-index: 5;
    }

    .btn-following {
        position: relative;
        z-index: 1;
        overflow: hidden;
        display: flex;
        align-items: flex-end;
    }

    .btn-following h4 {
        position: absolute;
        z-index: 10;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        color: var(--font-color, white);
    }

    .btn-following span {
        background-color: var(--background2-color, #111);
        height: 00px;
        width: 10%;
        position: relative;
        bottom: 0;
        z-index: 1;
    }

    .btn-following span:nth-child(1) {
        transition: all 1300ms ease;
    }
    .btn-following span:nth-child(2) {
        transition: all 1100ms ease;
    }
    .btn-following span:nth-child(3) {
        transition: all 900ms ease;
    }
    .btn-following span:nth-child(4) {
        transition: all 700ms ease;
    }
    .btn-following span:nth-child(5) {
        transition: all 500ms ease;
    }
    .btn-following span:nth-child(6) {
        transition: all 300ms ease;
    }
    .btn-following span:nth-child(7) {
        transition: all 500ms ease;
    }
    .btn-following span:nth-child(8) {
        transition: all 700ms ease;
    }
    .btn-following span:nth-child(9) {
        transition: all 900ms ease;
    }
    .btn-following span:nth-child(10) {
        transition: all 1100ms ease;
    }
    .btn-following span:nth-child(11) {
        transition: all 1300ms ease;
    }

    .btn-following:hover span {
        height: 122%;
    }

    .btn-following:active span {
        transition-duration: 0.1s;
        background-color: var(--background3-color, #222);
    }

    .btn-following:hover h4 {
        color: var(--font-hover-color, white);
    }

    .btn-following:hover {
        border-color: var(--border-hover-color, white);
    }

    .btn-following:active {
        border-color: var(--border-color, black);
    }
</style>