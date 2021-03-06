mport java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Graphics;
import java.awt.GridLayout;
import java.util.Arrays;

import javax.swing.JFrame;
import javax.swing.JLabel;
import javax.swing.JPanel;
import javax.swing.SwingUtilities;

public class ColorMapsTest
{
    public static void main(String[] args)
    {
        SwingUtilities.invokeLater(new Runnable()
        {
            @Override
            public void run()
            {
                createAndShowGUI();
            }
        });
    }

    private static void createAndShowGUI()
    {
        JFrame f = new JFrame();
        f.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);

        f.getContentPane().setLayout(new GridLayout(0,1));

        int steps = 1024;
        f.getContentPane().add(
            createPanel(steps, Color.RED));
        f.getContentPane().add(
            createPanel(steps, Color.RED, Color.GREEN));
        f.getContentPane().add(
            createPanel(steps, Color.RED, Color.GREEN, Color.BLUE));
        f.getContentPane().add(
            createPanel(steps,
                Color.RED, Color.YELLOW,
                Color.GREEN, Color.CYAN,
                Color.BLUE, Color.MAGENTA));
        f.getContentPane().add(
            createPanel(steps,
                Color.BLACK, Color.ORANGE, Color.WHITE,
                Color.BLUE, new Color(0,0,128)));


        JPanel panel = new JPanel(new BorderLayout());
        Color colors[] = new Color[]{ Color.RED, Color.GREEN, Color.BLUE };
        String info = "With sine over "+createString(colors);
        panel.add(new JLabel(info), BorderLayout.NORTH);
        ColorMapPanel1D colorMapPanel =
            new ColorMapPanel1D(
                ColorMaps1D.createSine(
                    ColorMaps1D.createDefault(steps, colors), Math.PI * 4));
        panel.add(colorMapPanel, BorderLayout.CENTER);
        f.getContentPane().add(panel);


        f.setSize(500, 400);
        f.setLocationRelativeTo(null);
        f.setVisible(true);
    }



    private static JPanel createPanel(int steps, Color ... colors)
    {
        JPanel panel = new JPanel(new BorderLayout());
        String info = "In "+steps+" steps over "+createString(colors);
        panel.add(new JLabel(info), BorderLayout.NORTH);
        ColorMapPanel1D colorMapPanel =
            new ColorMapPanel1D(ColorMaps1D.createDefault(steps, colors));
        panel.add(colorMapPanel, BorderLayout.CENTER);
        return panel;
    }

    private static String createString(Color ... colors)
    {
        StringBuilder sb = new StringBuilder();
        for (int i=0; i<colors.length; i++)
        {
            sb.append(createString(colors[i]));
            if (i < colors.length - 1)
            {
                sb.append(", ");
            }
        }
        return sb.toString();
    }
    private static String createString(Color color)
    {
        return "("+color.getRed()+","+color.getGreen()+","+color.getBlue()+")";
    }
}

// NOTE: This is an interface that is equivalent to the functional
// interface in Java 8. In an environment where Java 8 is available,
// this interface may be omitted, and the Java 8 version of this
// interface may be used instead.
interface DoubleFunction<R>
{
    R apply(double value);
}



/**
 * Interface for classes that can map a single value from the range
 * [0,1] to an int that represents an RGB color
 */
interface ColorMap1D
{
    /**
     * Returns an int representing the RGB color, for the given value in [0,1]
     *
     * @param value The value in [0,1]
     * @return The RGB color
     */
    int getColor(double value);
}

/**
 * Default implementation of a {@link ColorMap1D} that is backed by
 * a simple int array
 */
class DefaultColorMap1D implements ColorMap1D
{
    /**
     * The backing array containing the RGB colors
     */
    private final int colorMapArray[];

    /**
     * Creates a color map that is backed by the given array
     *
     * @param colorMapArray The array containing RGB colors
     */
    DefaultColorMap1D(int colorMapArray[])
    {
        this.colorMapArray = colorMapArray;
    }

    @Override
    public int getColor(double value)
    {
        double d = Math.max(0.0, Math.min(1.0, value));
        int i = (int)(d * (colorMapArray.length - 1));
        return colorMapArray[i];
    }
}


/**
 * Methods to create {@link ColorMap1D} instances
 */
class ColorMaps1D
{
    /**
     * Creates a {@link ColorMap1D} that walks through the given delegate
     * color map using a sine function with the given frequency
     *
     * @param delegate The delegate
     * @param frequency The frequency
     * @return The new {@link ColorMap1D}
     */
    static ColorMap1D createSine(ColorMap1D delegate, final double frequency)
    {
        return create(delegate, new DoubleFunction<Double>()
        {
            @Override
            public Double apply(double value)
            {
                return 0.5 + 0.5 * Math.sin(value * frequency);
            }
        });
    }

    /**
     * Creates a {@link ColorMap1D} that will convert the argument
     * with the given function before it is looking up the color
     * in the given delegate
     *
     * @param delegate The delegate {@link ColorMap1D}
     * @param function The function for converting the argument
     * @return The new {@link ColorMap1D}
     */
    static ColorMap1D create(
        final ColorMap1D delegate, final DoubleFunction<Double> function)
    {
        return new ColorMap1D()
        {
            @Override
            public int getColor(double value)
            {
                return delegate.getColor(function.apply(value));
            }
        };
    }


    /**
     * Creates a new ColorMap1D that maps a value between 0.0 and 1.0
     * (inclusive) to the specified color range, internally using the
     * given number of steps for interpolating between the colors
     *
     * @param steps The number of interpolation steps
     * @param colors The colors
     * @return The color map
     */
    static ColorMap1D createDefault(int steps, Color ... colors)
    {
        return new DefaultColorMap1D(initColorMap(steps, colors));
    }

    /**
     * Creates the color array which contains RGB colors as integers,
     * interpolated through the given colors.
     *
     * @param steps The number of interpolation steps, and the size
     * of the resulting array
     * @param colors The colors for the array
     * @return The color array
     */
    static int[] initColorMap(int steps, Color ... colors)
    {
        int colorMap[] = new int[steps];
        if (colors.length == 1)
        {
            Arrays.fill(colorMap, colors[0].getRGB());
            return colorMap;
        }
        double colorDelta = 1.0 / (colors.length - 1);
        for (int i=0; i<steps; i++)
        {
            double globalRel = (double)i / (steps - 1);
            int index0 = (int)(globalRel / colorDelta);
            int index1 = Math.min(colors.length-1, index0 + 1);
            double localRel = (globalRel - index0 * colorDelta) / colorDelta;

            Color c0 = colors[index0];
            int r0 = c0.getRed();
            int g0 = c0.getGreen();
            int b0 = c0.getBlue();
            int a0 = c0.getAlpha();

            Color c1 = colors[index1];
            int r1 = c1.getRed();
            int g1 = c1.getGreen();
            int b1 = c1.getBlue();
            int a1 = c1.getAlpha();

            int dr = r1-r0;
            int dg = g1-g0;
            int db = b1-b0;
            int da = a1-a0;

            int r = (int)(r0 + localRel * dr);
            int g = (int)(g0 + localRel * dg);
            int b = (int)(b0 + localRel * db);
            int a = (int)(a0 + localRel * da);
            int rgb =
                (a << 24) |
                (r << 16) |
                (g <<  8) |
                (b <<  0);
            colorMap[i] = rgb;
        }
        return colorMap;
    }

    /**
     * Private constructor to prevent instantiation
     */
    private ColorMaps1D()
    {
        // Private constructor to prevent instantiation
    }
}


/**
 * A panel painting a {@link ColorMap1D}
 */
class ColorMapPanel1D extends JPanel
{
    /**
     * The {@link ColorMap1D} that is painted
     */
    private final ColorMap1D colorMap;

    /**
     * Creates a new panel that paints the given color map
     *
     * @param colorMap The {@link ColorMap1D} to be painted
     */
    ColorMapPanel1D(ColorMap1D colorMap)
    {
        this.colorMap = colorMap;
    }

    @Override
    protected void paintComponent(Graphics g)
    {
        super.paintComponent(g);

        for (int x=0; x<getWidth(); x++)
        {
            double d = (double)x / (getWidth() - 1);
            int rgb = colorMap.getColor(d);
            g.setColor(new Color(rgb));
            g.drawLine(x, 0, x, getHeight());
        }

    }
}