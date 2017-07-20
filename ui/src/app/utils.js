
// Given a wheel event return if the event is upward or not.
function isWheelUp(e) {
    // http://www.javascriptkit.com/javatutors/onmousewheel.shtml
    let evt = window.event || e; //equalize event object
    if (e.nativeEvent !== undefined) {
        // Chromium
        evt = e.nativeEvent;
    }

    if (evt.detail == 0 && e.deltaY !== 0) {
        // https://developer.mozilla.org/en-US/docs/Web/Events/wheel
        return e.deltaY < 0;
    }
    return evt.detail? evt.detail < 0: evt.wheelDelta < 0; // make Opera use detail instead of wheelDelta;
}

export {isWheelUp};