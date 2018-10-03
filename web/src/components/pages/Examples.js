import React, { Component } from 'react'

class Examples extends Component {
  componentDidMount(){
    document.title = 'Examples - WPdirectory'
  }

  render() {
    return (
      <div className="page page-examples grid-container">
        <div className="grid-x grid-margin-x grid-margin-y">
          <div className="panel cell small-12">
            <h2>Examples</h2>
            <p>Writing regular expressions can be difficult so this page seeks to help you use WPDirectory by providing some examples and explaining what they do.</p>
            
            <h3>Functions</h3>
            <p>Let us start with something easy, let us search for use of the privacy function <code>wp_privacy_anonymize_data</code>, we can do this by searching for:</p>
            <pre>
                wp_privacy_anonymize_data\(
            </pre>
            <p>This works by matching the full name of the function when it is immediately followed by an opening parentheses <code>(</code>. Because parentheses <code>()</code> have a special meaning in regex we need to escape it <code>\(</code> to say that we want to match that exact character.</p>

            <h3>Hooks</h3>
            <p>Next up are hooks, let us search for the use of any hooks (actions/filters) beginning with <code>wp_privacy</code>, we can do this with:</p>
            <pre>
                add_(action|filter)\([\s|'|"]+wp_privacy\w*['|"]
            </pre>
            <p>This is a bit more complex so lets break it into sections.</p>
            <p>First we match adding actions and filters with <code>add_(action|filter)\(</code>, this matches both the add functions (add_action and add_filter) ending with the exact character match used above.</p>
            <p>Then we want to look for the start of the action/hook name but we cannot be sure what will come next, there may or may not be a space followed by single or double quotes. Using <code>[\s|'|"]+</code> lets us match one or more of the following- spaces, single quotes and double quotes.</p>
            <p>Now the main part, we want to match any hook name which begins with wp_privacy. We use <code>wp_privacy\w*</code> which matches our prefix and then any word characters (which includes underscores).</p>
            <p>Finally we finish off with <code>['|"]</code> which we use to ensure we reached the end of the hook name, as it must be followed by a single or double quote.</p>

          </div>
        </div>
      </div>
    )
  }
}

export default Examples