# Ray Tracing: The Next Week, in Go

This is a Go implementation of the raytracer described in the book [_Ray Tracing: The Next Week_](https://raytracing.github.io/books/RayTracingTheNextWeek.html) by Peter Shirley, Trevor David Black, Steve Hollasch.

It can render all images described in the book, that's why you may find some duplicated code. However, all the code that came from the previous project (Ray Tracing in One Weekend) has been cleaned up.

To generate an image run:

> go run . [image_number]

where __image_number__ is a number between 1 and 23.

Here's image #21, the famous Cornell Box, rendered with more than 33 billion rays:

![The Cornell Box, rendered with 210k rays per pixel](https://ascottix.github.io/rttnw/rttnw_cornell_box.png)

And the final image, showing off most of the features:

![Final image, rendered with 70k rays per pixel](https://ascottix.github.io/rttnw/rttnw_final.png)

Output is a file named `out.ppm` in PPM format.

All images are rendered with default parameter values. Different values can only be set by editing the source code.

## Note

To generate some images the file `earthmap.jpg` must be available in the project directory. It can be downloaded directly from the book page (it's image #4).
