import './App.css';
import React from "react";
import Webcam from "react-webcam";
import axios from 'axios';
import { ReactComponent as Camera } from "./camera.svg";
import Collapse from 'react-smooth-collapse';

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
        url: "http://localhost:8000",
        data: bodyFormData,
        headers: { "Content-Type": "multipart/form-data" },
      }).then((res) => {
        setRes(old => [{ data: res.data, src: ImageURL }, ...old]);
      }).catch(alert)
    },
    [webcamRef]
  );

  React.useEffect(() => {
    window.addEventListener('keyup', capture);
    return () => {
      window.removeEventListener('keyup', capture);
    };
  }, [capture]);

  return (
    <main>
      <div className='camera-container'>
        <Webcam
          audio={false}
          width={1280}
          height={720}
          ref={webcamRef}
          screenshotFormat="image/jpeg"
          videoConstraints={videoConstraints}
          mirrored={true}
        />
        <button className='button' onClick={capture}><Camera style={{
          height: '100%',
          width: '100%',
          fill: '#666',
        }} /></button>
      </div>
      <div className='showroom'>
        {res.map((r, i) => <Pic key={i} src={r.src} data={r.data} />)}
      </div>
    </main>
  );
};

const Pic = ({ key, src, data }) => {
  const [open, change] = React.useState(false);

  return (
    <div key={key} >
      <img src={src} alt="img"
        style={ Object.keys(data).length === 0 ?
          { background: '#d337' } :
          { background: '#9d67' }}
        onClick={() => {
          change((open) => !open);
        }} />
      <Collapse expanded={open}>
        <pre>{JSON.stringify(data, "", "  ")}</pre>
      </Collapse>
    </div>
  )
}

export default App;