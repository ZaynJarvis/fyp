import './App.css';
import React from "react";
import Webcam from "react-webcam";
import axios from 'axios';

const videoConstraints = {
  width: 1280,
  height: 720,
  facingMode: "user"
};

function b64toBlob(b64Data, contentType, sliceSize) {
  contentType = contentType || '';
  sliceSize = sliceSize || 512;

  var byteCharacters = atob(b64Data);
  var byteArrays = [];

  for (var offset = 0; offset < byteCharacters.length; offset += sliceSize) {
    var slice = byteCharacters.slice(offset, offset + sliceSize);

    var byteNumbers = new Array(slice.length);
    for (var i = 0; i < slice.length; i++) {
      byteNumbers[i] = slice.charCodeAt(i);
    }

    var byteArray = new Uint8Array(byteNumbers);

    byteArrays.push(byteArray);
  }

  var blob = new Blob(byteArrays, { type: contentType });
  return blob;
}

const App = () => {
  const webcamRef = React.useRef(null);
  const [res, setRes] = React.useState([]);

  const capture = React.useCallback(
    () => {
      const ImageURL = webcamRef.current.getScreenshot();

      const bodyFormData = new FormData();
      const block = ImageURL.split(";");
      const contentType = block[0].split(":")[1];
      const realData = block[1].split(",")[1];
      const blob = b64toBlob(realData, contentType);

      bodyFormData.append("image", blob);
      axios({
        method: "post",
        url: "http://0.0.0.0:8000",
        data: bodyFormData,
        headers: { "Content-Type": "multipart/form-data" },
      }).then((res) => {
        setRes(old => [{ data: res.data, src: ImageURL }, ...old]);
      }).catch(alert)
    },
    [webcamRef]
  );

  return (
    <>
      <Webcam
        audio={false}
        width={640}
        height={360}
        ref={webcamRef}
        screenshotFormat="image/jpeg"
        videoConstraints={videoConstraints}
        mirrored={true}
      />
      <button style={{
        position: 'fixed',
        top: '1rem',
        left: '1rem',
        borderRadius: '3px',
        fontSize: '2rem',
      }} onClick={capture}>Capture</button>
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(5, 1fr)',
      }}>
        {res.map((r, i) => (<div key={i}>
          <img style={{ height: '180px', width: '320px' }} src={r.src} alt="img" />
          <pre>{JSON.stringify(r.data, "", "  ")}</pre>
        </div>))}
      </div>
    </>
  );
};

export default App;