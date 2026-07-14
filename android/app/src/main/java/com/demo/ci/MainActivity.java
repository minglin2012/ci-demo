package com.demo.ci;

import android.os.Bundle;
import android.widget.TextView;
import androidx.appcompat.app.AppCompatActivity;

/**
 * Main Activity for CI Demo Android App.
 * Displays build information to verify the CI pipeline.
 */
public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        TextView textView = new TextView(this);
        textView.setText(
            "╔══════════════════════════════════════╗\n" +
            "║   CI Demo Android App v1.0          ║\n" +
            "╚══════════════════════════════════════╝\n\n" +
            "✅ Built via GitHub Actions CI/CD\n" +
            "🤖 Android Platform Build\n" +
            "📱 Min SDK: 24 | Target SDK: 34\n" +
            "\n" +
            "This APK was automatically assembled\n" +
            "by the multi-platform CI pipeline."
        );
        textView.setPadding(48, 48, 48, 48);
        textView.setTextSize(16);
        setContentView(textView);
    }
}
