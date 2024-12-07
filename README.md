# Shadfin >> Desktop

**Windows only for now.**

You can either install via the [artifact in the workflow](https://github.com/Shadfin/desktop/actions), or by the [latest release published](https://github.com/Shadfin/desktop/releases).

Additional documentation will be provided in the future.


# How it Works

Shadfin takes a different approach from typical apps or players that use MPV.

Instead of having to use the same OpenGL context/buffer for both the webview and the video, we can use two completely separate contexts, or buffers. This eliminates issue pertaining to have and sync up the UI with the refresh rate of the content.

## The Issue: Syncing the Frame

This has a major disadvantage, we have to ensure the webview and the player are synced up for one to render after the other.

Typically this means we are stuck with syncing each frame of the webview with each frame of the player. 
However doing so causes the webview to feel noticeably slower.

A video at 24fps will cause the webview to also run at 24 fps, not good! 

![frame-copy](https://github.com/user-attachments/assets/e671fd11-f826-49ca-8461-6b5f96eb96f9)

### Solution 1: Don't sync the frames?

While it may seem obvious this solution is not always guaranteed to work.<br>
We are no longer ensuring the webview waits until the video frame has been rendered, this removes the guarantee we had before of the webview always rendering on top of the video content.

This not only disrupts the ability to use the webview it causes a noticeable flicker as the webview and player try to compete who gets to be on top.

If you've ever heard of "Z-Fighting" in game development this is essentially the same exact thing. We are wanting two things to be rendered at the same depth (or the Z cord) and they will fight each other to be on the top.

![2017-10-04_22-53-10](https://github.com/user-attachments/assets/175c14cc-8cd4-4e4b-876c-994e3759e6f1)

## Screw Sync - Maintain Order

Turns out, we can have our cake and eat it too.

**Why don't we just render two windows?**

Parent Window: Holds the actual player, and titlebar.

Child Window: Holds the webview.


Turns out this approach greatly reduces the complexity as well, sweet!

In Windows at least a parent window can have multiple child windows, that are rendered on top of the parent. This fixes the "Z-Fighting" explained before, we have the underlying window manager figure it out for us, which is *much* faster than doing it ourselves. 

And boom! A crazy fast native player with web controls!

#### Whats the catch?

Turns out windows is rather unique with supporting this, on MacOS theoretically this should also be supported with minimal hassle, but I've yet to test it.

On linux the story gets even more complicated. X11 being so old doesn't really support a use case like this, and Wayland doesn't really have the concepts of 'windows' anymore, so we're probably stuck falling back to syncing the frames. But who knows! Maybe we can get Linux to use two windows too!
