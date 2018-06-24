import React, { Component } from 'react'

class Logo extends Component {  
  render() {
    
    const {
        width = 1040,
        height = 200,
    } = this.props;

    return (
        <svg
            aria-hidden
            role="img"
            focusable="false"
            className="logo"
            xmlns="http://www.w3.org/2000/svg"
            width={ width }
            height={ height }
            viewBox="0 0 1043.949 200.99"
            enable-background="0 0 1043.949 200.99"
        ><g><g>
            <path fill="#39B54A" d="M48.953,67.134L32.146,50.321C12.109,80.66,12.109,120.325,32.14,150.656l16.813-16.806 C37.322,113.203,37.322,87.789,48.953,67.134z"/>
            <path fill="#F7941E" d="M74.905,159.81l-16.813,16.806c30.332,20.037,70.003,20.037,100.342,0l-16.806-16.806 C120.974,171.44,95.56,171.44,74.905,159.81z"/>
            <path fill="#00AEEF" d="M141.621,41.182l16.807-16.813c-30.333-20.03-69.997-20.03-100.329,0l16.807,16.813 C95.56,29.552,120.974,29.552,141.621,41.182z"/>
            <path fill="#D24632" d="M184.387,50.328L167.58,67.134c11.624,20.647,11.631,46.062,0,66.716l16.807,16.806 C204.423,120.325,204.423,80.66,184.387,50.328z"/>
            <path fill="#4A4A4A" d="M108.267,156.996c31.2,0,56.5-25.292,56.5-56.5s-25.299-56.507-56.5-56.507 c-31.208,0-56.507,25.3-56.507,56.507c0,11.063,3.195,21.379,8.687,30.095l-24.875,24.876c5.068,6.69,11.034,12.656,17.726,17.718 l24.875-24.876C86.887,153.808,97.196,156.996,108.267,156.996z M64.984,100.496c0-23.906,19.376-43.283,43.283-43.283 c23.899,0,43.276,19.376,43.276,43.283c0,23.899-19.376,43.275-43.276,43.275C84.36,143.771,64.984,124.396,64.984,100.496z"/>
        </g><g>
            <path fill="#4A4A4A" d="M240.732,147.011l-22.086-93.04h22.504l7.034,38.381c2.078,11.043,4.009,23.051,5.521,32.434h0.278 c1.521-10.078,3.73-21.26,6.078-32.712l7.869-38.103h22.355l7.452,39.207c2.078,10.904,3.591,20.842,4.974,31.19h0.278 c1.374-10.348,3.452-21.251,5.382-32.294l7.591-38.103h21.398l-24.024,93.04h-22.773l-7.869-40.033 c-1.791-9.382-3.313-18.078-4.417-28.712H278c-1.652,10.495-3.174,19.33-5.383,28.712l-8.834,40.033H240.732z"/>
            <path fill="#4A4A4A" d="M348.285,55.214c6.486-1.104,15.599-1.93,28.434-1.93c12.973,0,22.225,2.487,28.434,7.452 c5.938,4.696,9.938,12.426,9.938,21.539c0,9.104-3.035,16.834-8.556,22.086c-7.182,6.756-17.808,9.8-30.233,9.8 c-2.756,0-5.244-0.14-7.174-0.417v33.268h-20.842V55.214z M369.127,97.456c1.792,0.409,4,0.548,7.043,0.548 c11.173,0,18.077-5.652,18.077-15.182c0-8.557-5.938-13.669-16.425-13.669c-4.278,0-7.182,0.417-8.695,0.834V97.456z"/>
            <path fill="#4A4A4A" d="M456.384,55.214c7.73-1.243,17.808-1.93,28.434-1.93c17.669,0,29.129,3.174,38.094,9.938 c9.669,7.174,15.739,18.634,15.739,35.06c0,17.808-6.487,30.095-15.46,37.686c-9.799,8.147-24.704,12.008-42.929,12.008 c-10.904,0-18.634-0.687-23.877-1.383V55.214z M477.505,131.133c1.791,0.417,4.687,0.417,7.313,0.417 c19.052,0.139,31.469-10.356,31.469-32.573c0.139-19.33-11.174-29.547-29.26-29.547c-4.695,0-7.73,0.417-9.521,0.835V131.133z"/>
            <path fill="#4A4A4A" d="M574.014,60.736c0,5.8-4.417,10.496-11.321,10.496c-6.625,0-11.043-4.696-10.903-10.496 c-0.14-6.069,4.278-10.625,11.043-10.625C569.597,50.11,573.875,54.667,574.014,60.736z M552.337,147.011V79.509h20.981v67.501 H552.337z"/>
            <path fill="#4A4A4A" d="M590.162,101.734c0-9.938-0.27-16.425-0.548-22.225h18.078l0.696,12.425h0.548 c3.452-9.799,11.738-13.947,18.225-13.947c1.93,0,2.896,0,4.417,0.278v19.738c-1.522-0.27-3.313-0.548-5.661-0.548 c-7.73,0-12.973,4.139-14.355,10.626c-0.278,1.382-0.418,3.035-0.418,4.695v34.233h-20.981V101.734z"/>
            <path fill="#4A4A4A" d="M657.671,119.951c0.688,8.704,9.244,12.843,19.043,12.843c7.183,0,12.982-0.965,18.643-2.765l2.757,14.217 c-6.904,2.765-15.321,4.148-24.434,4.148c-22.912,0-36.024-13.252-36.024-34.373c0-17.121,10.625-36.033,34.094-36.033 c21.808,0,30.095,16.982,30.095,33.686c0,3.591-0.417,6.765-0.696,8.278H657.671z M682.375,105.604 c0-5.113-2.209-13.669-11.869-13.669c-8.835,0-12.426,8.009-12.974,13.669H682.375z"/>
            <path fill="#4A4A4A" d="M765.076,145.358c-3.721,1.652-10.764,3.035-18.764,3.035c-21.816,0-35.756-13.252-35.756-34.512 c0-19.738,13.53-35.895,38.651-35.895c5.521,0,11.591,0.974,16.008,2.626l-3.313,15.6c-2.479-1.104-6.209-2.069-11.73-2.069 c-11.043,0-18.226,7.869-18.086,18.912c0,12.426,8.286,18.912,18.503,18.912c4.966,0,8.835-0.834,12.009-2.069L765.076,145.358z" />
            <path fill="#4A4A4A" d="M802.353,60.188v19.321h15.052v15.46h-15.052v24.434c0,8.147,1.931,11.869,8.287,11.869 c2.617,0,4.687-0.278,6.208-0.548l0.14,15.869c-2.766,1.104-7.73,1.8-13.67,1.8c-6.765,0-12.416-2.348-15.738-5.8 c-3.86-4-5.791-10.487-5.791-20.017V94.969h-8.974v-15.46h8.974V64.875L802.353,60.188z"/>
            <path fill="#4A4A4A" d="M895.541,112.499c0,24.712-17.529,36.033-35.616,36.033c-19.738,0-34.92-12.982-34.92-34.79 c0-21.808,14.356-35.755,36.024-35.755C881.732,77.987,895.541,92.213,895.541,112.499z M846.673,113.195 c0,11.591,4.834,20.286,13.808,20.286c8.139,0,13.391-8.139,13.391-20.286c0-10.078-3.869-20.295-13.391-20.295 C850.404,92.9,846.673,103.256,846.673,113.195z"/>
            <path fill="#4A4A4A" d="M908.523,101.734c0-9.938-0.27-16.425-0.548-22.225h18.086l0.687,12.425h0.549 c3.451-9.799,11.738-13.947,18.225-13.947c1.931,0,2.896,0,4.418,0.278v19.738c-1.522-0.27-3.313-0.548-5.661-0.548 c-7.73,0-12.974,4.139-14.356,10.626c-0.277,1.382-0.417,3.035-0.417,4.695v34.233h-20.981V101.734z"/>
            <path fill="#4A4A4A" d="M979.207,79.509l10.078,33.129c1.104,4.008,2.486,8.974,3.313,12.565h0.418 c0.965-3.591,2.069-8.695,3.034-12.565l8.278-33.129h22.503l-15.738,44.45c-9.66,26.782-16.147,37.546-23.738,44.311 c-7.321,6.348-15.052,8.557-20.295,9.243l-4.417-17.808c2.626-0.409,5.939-1.652,9.113-3.583c3.174-1.661,6.625-4.974,8.695-8.426 c0.687-0.965,1.104-2.07,1.104-3.035c0-0.687-0.139-1.791-0.965-3.452l-24.712-61.702H979.207z"/>
        </g></g></svg>
		);
  }
}

export default Logo